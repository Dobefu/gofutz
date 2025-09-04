interface Line {
  number: number;
  startLine: number;
  endLine: number;
  startColumn: number;
  endColumn: number;
  numberOfStatements: number;
  executionCount: number;
}

interface TestStatus {
  Pending: 0;
  Running: 1;
  Passed: 2;
  Failed: 3;
  NoTests: 4;
}

interface File {
  name: string;
  functions: Function[];
  code: string;
  highlightedCode: string;
  status: TestStatus;
  coverage: number;
  coveredLines: Line[];
}

interface Function {
  name: string;
  result: TestResult;
}

interface TestResult {
  coverage: number;
}

interface Params {
  files: Record<string, File>;
}

interface Message {
  method: string;
  error: string;
  params: Params;
}
