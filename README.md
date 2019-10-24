gentest 用于快速生成基于 gomock + testify 风格的单元测试

执行gentest new 查看参数说明

example after creating unit test with gentest new ：

func TestCheck(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock data here
	mockEngine := mock.NewMockEngine(ctrl)
	mockEngine.EXPECT().Run().Return()

	// call func here after mock

	a.Equal(nil, nil)

}