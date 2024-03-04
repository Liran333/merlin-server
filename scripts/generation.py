import os
import base64
import hashlib
import argparse
from getpass import getpass
from binascii import hexlify
from typing import Tuple

from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.backends import default_backend


def encrypt_token(token: str) -> Tuple[str, str]:
    salt = os.urandom(32)
    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA256(),
        length=32,
        salt=salt,
        iterations=10000,
        backend=default_backend(),
    )
    key = kdf.derive(token.encode())
    return base64.b64encode(key).decode(), base64.b64encode(salt).decode()


def encode_token(token: str, salt: str) -> str:
    salt_bytes = base64.b64decode(salt)
    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA256(),
        length=32,
        salt=salt_bytes,
        iterations=10000,
        backend=default_backend(),
    )
    key = kdf.derive(token.encode())
    return base64.b64encode(key).decode()


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("token", help="the token to be encrypted")
    args = parser.parse_args()

    enc, salt = encrypt_token(args.token)

    encLog = enc.replace("=", "")
    saltLog = salt.replace("=", "")

    print(f"enc: {encLog}")
    print(f"salt: {saltLog}")

    # Write the salt to the .env file
    with open(".env", "r") as file:
        data = file.readlines()

    for i, line in enumerate(data):
        if line.startswith("INTERNAL_SALT="):
            data[i] = f"INTERNAL_SALT={saltLog}\n"
        if line.startswith("INTERNAL_TOKEN_HASH="):
            data[i] = f"INTERNAL_TOKEN_HASH={encLog}\n"

    with open(".env", "w") as file:
        file.writelines(data)


if __name__ == "__main__":
    main()
