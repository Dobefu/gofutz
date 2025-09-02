interface File {
  name: string;
  tests: Test[];
  code: string;
  highlightedCode: string;
}

interface Test {
  name: string;
  result: TestResult;
}

interface TestStatus {
  Pending: 0;
  Running: 1;
  Passed: 2;
  Failed: 3;
}

interface Line {
  number: number;
  executionCount: number;
}

interface TestResult {
  status: TestStatus;
  output: string[];
  coverage: number;
  coveredLines: Line[];
}

interface Params {
  files: Record<string, File>;
}

interface Message {
  method: string;
  error: string;
  params: Params;
}
