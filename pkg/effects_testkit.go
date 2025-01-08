package pkg

type TestKitArgs struct {
	Wd       string
	Files    map[string][]byte
	QaBool   map[string]bool
	QaString map[string]string
}

//func NewEffectsTestkit () *TestKit {
//
//}

//type TestKit struct {
//	MockLogger *MockLogger
//	MockFs     *MockFs
//	MockExec   *MockExec
//	MockOs     *MockOs
//	Effects         *SystemUtils
//	cfg        *Config
//	cmd        *Commands
//}
//
//// NewTestKit creates a new TestKit
//// wd: working directory from where the command is executed
//// files: map of files with their content (pass nil if not needed)
//// qABool: map of questions and answers for boolean questions (pass nil if not needed)
//// qaString: map of questions and answers for string questions (pass nil if not needed)
//func NewTestKit(wd string, files map[string][]byte, qABool map[string]bool, qaString map[string]string) (tk *TestKit, err error) {
//	mockFs := NewMockFs(files)
//	mockExec := NewMockExec()
//	mockLogger := NewMockLogger()
//	mockOs := NewMockOs(wd, qABool, qaString)
//	su := NewSystemUtils(mockFs, mockExec, mockLogger, mockOs)
//	cfg, err := NewConfig(su)
//	if err != nil {
//		return &TestKit{}, err
//	}
//	return &TestKit{
//		MockLogger: mockLogger,
//		MockFs:     mockFs,
//		MockExec:   mockExec,
//		MockOs:     mockOs,
//		su:         su,
//		cfg:        cfg,
//		cmd:        NewCommands(su, cfg),
//	}, nil
//}
