package giteaimpl

import "testing"

var impl giteaImpl

func TestParseTags(t *testing.T) {
	readMeFile := map[string]string{
		"empty":            "",
		"frontMatterEmpty": "xxxxxxxxxxfdsf",
		"normal":           "---\ninference: false\nlicense: other\nlicense_name: microsoft-research-license\nlicense_link: https://huggingface.co/microsoft/phi-2/resolve/main/LICENSE\nlanguage:\n- en\npipeline_tag: text-generation\ntags:\n- nlp\n- code\nframeworks:\n- xxxxx\n- yyyyy\ntask: xxxx\n---\n\n## Model Summary\n\nPhi-2 is a Transformer with **2.7 billion** parameters. It was trained using the same data sources as [Phi-1.5](https://huggingface.co/microsoft/phi-1.5), augmented with a new data source that consists of various NLP synthetic texts and filtered websites (for safety and educational value). When assessed against benchmarks testing common sense, language understanding, and logical reasoning, Phi-2 showcased a nearly state-of-the-art performance among models with less than 13 billion parameters.\n\nOur model hasn't been fine-tuned through reinforcement learning from human feedback. The intention behind crafting this open-source model is to provide the research community with a non-restricted small model to explore vital safety challenges, such as reducing toxicity, understanding societal biases, enhancing controllability, and more.",
	}

	data, err := impl.parseTags([]byte(readMeFile["empty"]))
	if err == nil {
		t.Fatal("case empty fail")
	}

	data, err = impl.parseTags([]byte(readMeFile["frontMatterEmpty"]))
	if err == nil {
		t.Fatal("case frontMatterEmpty fail")
	}

	data, err = impl.parseTags([]byte(readMeFile["normal"]))
	if err != nil {
		t.Fatalf("case normal fail: %s", err.Error())
	}

	if len(data.Tags) != 2 || len(data.Frameworks) == 0 || data.Task == "" {
		t.Fatal("case normal fail")
	}
}
