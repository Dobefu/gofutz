interface File {
  name: string;
  tests: Test[];
  code: string;
  highlightedCode: string;
}

interface Test {
  name: string;
}

interface Params {
  files: File[];
}

interface Message {
  method: string;
  error: string;
  params: Params;
}
