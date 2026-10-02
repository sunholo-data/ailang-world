package main
import("go/parser";"go/token";"go/ast";"fmt";"os";"strings")
func main(){b,_:=os.ReadFile("/tmp/world-iter221-product-eval-independent/changed-go.txt");for _,p:=range strings.Split(string(b),"\n"){s:=token.NewFileSet();f,e:=parser.ParseFile(s,p,nil,0);if e!=nil{panic(e)};for _,d:=range f.Decls{if fn,ok:=d.(*ast.FuncDecl);ok{n:=s.Position(fn.End()).Line-s.Position(fn.Pos()).Line+1;if n>50{fmt.Printf("%s:%d %s %d\n",p,s.Position(fn.Pos()).Line,fn.Name.Name,n)}}}}}
