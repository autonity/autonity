package graph

type Graph struct {
	//graph bool `@"g""r""a""p""h"" ""L""R"";"`
	Edges []*Edge `@@*`
}

type Edge struct {
	LeftNode  string `@Ident[" "|"\t"]"-""-"["-"]`
	Directed  bool   `[@">"][" "|"\t"]`
	RightNode string `@String[";"]`
}
