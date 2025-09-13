declare global {
  interface Window {
    documents: Record<string, File>;
  }
}

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

interface InitParams {
  files: Record<string, File>;
  coverage: number;
  isRunning: boolean;
  output: string[];
}

interface InitMessage {
  method: string;
  error: string;
  params: InitParams;
}

interface UpdateParams {
  files: Record<string, File>;
  coverage: number;
  isRunning: boolean;
}

interface UpdateMessage {
  method: string;
  error: string;
  params: UpdateParams;
}

interface OutputParams {
  output: string[];
}

interface OutputMessage {
  method: string;
  error: string;
  params: OutputParams;
}
