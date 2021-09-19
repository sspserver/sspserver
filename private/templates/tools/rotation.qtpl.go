// Code generated by qtc from "rotation.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/tools/rotation.qtpl:1
package tools

//line templates/tools/rotation.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/tools/rotation.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/tools/rotation.qtpl:1
func StreamRotationToolTpl(qw422016 *qt422016.Writer) {
//line templates/tools/rotation.qtpl:1
	qw422016.N().S(`<!DOCTYPE html><html><head><meta charset="UTF-8" /><title>Rotation tool</title><link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css" /><script src="https://unpkg.com/react@16/umd/react.development.js"></script><script src="https://unpkg.com/react-dom@16/umd/react-dom.development.js"></script><script src="https://unpkg.com/babel-standalone@6.15.0/babel.min.js"></script><script src="https://cdnjs.cloudflare.com/ajax/libs/axios/0.18.0/axios.min.js"></script></head><body><nav class="navbar navbar-expand-lg navbar-light bg-light"><a class="navbar-brand" href="#">Rotation</a><button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation"><span class="navbar-toggler-icon"></span></button><div class="collapse navbar-collapse" id="navbarSupportedContent"><!-- <ul class="navbar-nav mr-auto"><li class="nav-item active"><a class="nav-link" href="#">Home <span class="sr-only">(current)</span></a></li><li class="nav-item"><a class="nav-link" href="#">Link</a></li><li class="nav-item dropdown"><a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Dropdown</a><div class="dropdown-menu" aria-labelledby="navbarDropdown"><a class="dropdown-item" href="#">Action</a><a class="dropdown-item" href="#">Another action</a><div class="dropdown-divider"></div><a class="dropdown-item" href="#">Something else here</a></div></li><li class="nav-item"><a class="nav-link disabled" href="#">Disabled</a></li></ul> --><form class="form-inline my-2 my-lg-0"><input class="form-control mr-sm-2" type="search" placeholder="Search" aria-label="Search"><button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button></form></div></nav><div class="container"><div class="clear"><p>&nbsp;</p></div><div id="root"></div></div><script type="text/babel">class ListAdItemComponent extends React.Component {render() {console.log(this.props.items);return this.props.items.map((item, index) => <tr key={index}><td>{item.ID} [{item.Opt ? item.Opt[0] : -1}] {item.Format.codename}</td><td>{item.State.spent / 1000000000} / {item.DailyBudget / 1000000000}</td><td>{item.State.impressions} / {item.State.clicks} / {item.State.leads}</td><td>{item.Price / 1000000000} / {item.LeadPrice / 1000000000} / {item.BidPrice / 1000000000}</td></tr>)}}class ListComponent extends React.Component {constructor(props) {super(props);}componentDidMount() {let self = this;axios.get("/v1/tools/rotation.json").then(function(response) {self.setState({campaigns: response.data});})}render() {if (!this.state || !this.state.campaigns) {return <div>NOOO</div>}return (<div><table className="table table-hover"><thead className="thead-dark"><tr><th scope="col" style={{width:"3em"}}>#</th><th scope="col">#ID</th><th scope="col">Spent</th><th scope="col">Imps/Clicks/Leads</th><th scope="col">Price/Lead/BidPrice</th></tr></thead><tbody>{this.state.campaigns.map((item, index) => <React.Fragment key={index}><tr><td className="text-righ" rowSpan={item.Ads.length + 1}>{index + 1}</td><td>{item.ID}</td><td></td><td></td><td></td></tr><ListAdItemComponent items={item.Ads} /></React.Fragment>)}</tbody></table></div>);}toJSON(data) {return JSON.stringify(data);}}ReactDOM.render(<ListComponent />,document.getElementById('root'));</script><script src="https://code.jquery.com/jquery-3.3.1.slim.min.js"></script><script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"></script><script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/js/bootstrap.min.js"></script></body></html>`)
//line templates/tools/rotation.qtpl:131
}

//line templates/tools/rotation.qtpl:131
func WriteRotationToolTpl(qq422016 qtio422016.Writer) {
//line templates/tools/rotation.qtpl:131
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/tools/rotation.qtpl:131
	StreamRotationToolTpl(qw422016)
//line templates/tools/rotation.qtpl:131
	qt422016.ReleaseWriter(qw422016)
//line templates/tools/rotation.qtpl:131
}

//line templates/tools/rotation.qtpl:131
func RotationToolTpl() string {
//line templates/tools/rotation.qtpl:131
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/tools/rotation.qtpl:131
	WriteRotationToolTpl(qb422016)
//line templates/tools/rotation.qtpl:131
	qs422016 := string(qb422016.B)
//line templates/tools/rotation.qtpl:131
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/tools/rotation.qtpl:131
	return qs422016
//line templates/tools/rotation.qtpl:131
}
