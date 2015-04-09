package dashboard

const (
	tpl = `
<html>
<head>
<title>{{ .Title }}</title>
<meta http-equiv="refresh" content="{{ .Refresh}}">
<style>
#content {
	margin: 0 auto;
	padding: 10px;
}

.container {
	box-sizing: border-box;
	width: 1200px;
	height: 450px;
	padding: 20px 15px 15px 15px;
	margin: 15px auto 30px auto;
	border: 1px solid #ddd;
	background: #fff;
	background: linear-gradient(#f6f6f6 0, #fff 50px);
	background: -o-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -ms-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -moz-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -webkit-linear-gradient(#f6f6f6 0, #fff 50px);
	box-shadow: 0 3px 10px rgba(0,0,0,0.15);
	-o-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-ms-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-moz-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-webkit-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
}

.graph {
	width: 100%;
	height: 100%;
	font-size: 14px;
	line-height: 1.2em;
}
</style>

<script src="http://cdnjs.cloudflare.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
<script src="http://cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.min.js"></script>
<script src="http://cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.time.min.js"></script>
<script src="http://cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.selection.min.js"></script>
<script type="text/javascript">
	var options = {
		legend: {
			position: "nw",
			noColumns: 5,
			backgroundOpacity: 0.2
		},
		xaxis: {
			mode: "time",
			timezone: "browser",
			timeformat: "%H:%M:%S "
		},		
	};

  {{range $idx, $graph := .Graphs}}
	var data_{{$idx}} = [
	{{range $lidx, $line := $graph.Lines}}
		{label: "{{$line.Legend}}", data: {{$line.Points}}},
	{{end}}
	];
  {{end}}		

	$(document).ready(function() {
	{{range $idx, $graph := .Graphs}}
	  $.plot("#graph{{$idx}}", data_{{$idx}}, options);	
  	{{end}}		
	});	
</script>

</head>
<body>

<div id="content">
  {{range $idx, $graph := .Graphs}}
	<p>{{$graph.Title}}</p>
	<div class="container" style="height:200px;">		
		<div id="graph{{$idx}}" class="graph"></div>
	</div>	
  {{end}}
</div>

</body>
</html>
	`
)
