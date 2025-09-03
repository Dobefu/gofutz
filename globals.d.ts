interface Line {
  number: number;
  startLine: number;
  endLine: number;
  startColumn: number;
  endColumn: number;
  numberOfStatements: number;
  executionCount: number;
}

interface File {
  name: string;
  functions: Function[];
  code: string;
  highlightedCode: string;
  coverage: number;
  coveredLines: Line[];
}

interface Function {
  name: string;
  result: TestResult;
}

interface TestStatus {
  Pending: 0;
  Running: 1;
  Passed: 2;
  Failed: 3;
}

interface TestResult {
  status: TestStatus;
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
