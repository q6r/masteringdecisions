function main(body)
{
/*
    $('title')
        .html('title goes here');

    $('<p>')
        .html('Data : ' + JSON.stringify(body))
        .appendTo('body');
*/
	
	
	var userName = "Admin"
	var decisions_inProgress = ["decision1", "decision2", "decision3","decision4"]
	var decisions_completed = ["decision5", "decision6", "decision7","decision8"]
	
	//nav section
	var nav = $('<nav class="navbar navbar-inverse navbar-fixed-top">').appendTo('body')
	var div_container = $('<div class="container-fluid">').appendTo(nav)
	var div_nav_header = $('<div class="navbar-header">').appendTo(div_container)
						
					
	var button_nav = $([
			'<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#myNavbar">',
						'<span class="icon-bar"></span>',
						'<span class="icon-bar"></span>',
						'<span class="icon-bar"></span>',
					'</button>'
	].join('\n'));
					
	div_nav_header.append(button_nav)	
	div_nav_header.append($('<a class = "navbar-brand" href="#"><img id="logo" src="new-woodentable.png">'))
	
	var div_collapse = $('<div class="collapse navbar-collapse" id="myNavbar">')
	
	var nav_ul1 = $('<ul class="nav navbar-nav"><li> <a href="#">Dashboard</a></li></ul>')
	var nav_ul2 = $([
	'<ul class="nav navbar-nav navbar-right">',
		'<li class="dropdown">',
			'<a class="dropdown-toggle" role="button" data-toggle="dropdown" href="#" aria-expanded="false"><i class="glyphicon glyphicon-user"></i>' + userName+ '<span class="caret"></span></a>',
				'<ul id="g-account-menu" class="dropdown-menu" role="menu">',
					'<li><a href="#">My Profile</a></li>',
				'</ul>',
		'</li>',
		'<li><a href="#"><i class="glyphicon glyphicon-lock"></i> Logout</a></li>',
	'</ul>'
	].join("\n"));
				
				
	div_collapse.append(nav_ul1)
	div_collapse.append(nav_ul2)
	div_container.append(div_collapse)
	
	//dashboard section
	
	var div_dashboard = $('<div class="container-fluid" id="my-own-style">').appendTo('body')
	
	var div_row = $('<div class="row-fluid" >')
	
	var nav_section = $([
			
			'<div class="col-sm-3" >',
				'<a href ="#">',
					'<strong><i class="glyphicon glyphicon-wrench"></i> Tools</strong>',
				'</a>',
				'<hr>',
				'<ul class="nav nav-stacked">',
				'<li class="nav-header"> <a href="#" data-toggle="collapse" data-target="#userMenu" aria-expanded="false" class="collapsed">Decisions <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
                    '<ul class="nav nav-stacked collapse" id="userMenu" aria-expanded="false" style="height: 0px;">',
						
                        '<li class="active"> <a href="#"><i class="glyphicon glyphicon-asterisk"></i> New Decision</a></li>',
						'<li class="nav-header"> <a href="#" data-toggle="collapse" data-target="#userMenu3" aria-expanded="false" class="collapsed">In Progress <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
							'<ul class="nav nav-stacked collapse" id="userMenu3" aria-expanded="false" style="height: 0px;">',
								'<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Decision1 </a></li>',
								'<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Decision2 </a></li>',
								'<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Decision3 </a></li>',
								'<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Decision4 </a></li>',
							'</ul>',
						'</li>',
                       '<li class="nav-header"> <a href="#" data-toggle="collapse" data-target="#userMenu4" aria-expanded="false" class="collapsed">Completed <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
							'<ul class="nav nav-stacked collapse" id="userMenu4" aria-expanded="false" style="height: 0px;">',
								'<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Decision5 </a></li>',
								'<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Decision6 </a></li>',
								'<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Decision7 </a></li>',
								'<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Decision8 </a></li>',
							'</ul>',
						'</li>',
                    '</ul>',
                '</li>',
				'<li class="nav-header"> <a href="#" data-toggle="collapse" data-target="#Menu2" aria-expanded="false" class="collapsed">Settings <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
					
                    '<ul class="nav nav-stacked collapse" id="Menu2" aria-expanded="fasle">',
                        '<li class="active"> <a href="#"><i class="glyphicon glyphicon-home"></i> Home</a></li>',
                        
                        //'<li><a href="#"><i class="glyphicon glyphicon-cog"></i> Options</a></li>',
                        
                        '<li><a href="/static/person.html"><i class="glyphicon glyphicon-user"></i> New User</a></li>',
                        '<li><a href="#"><i class="glyphicon glyphicon-list-alt"></i> Report</a></li>',
                        
                        '<li><a href="#"><i class="glyphicon glyphicon-off"></i> Logout</a></li>',
						
                    '</ul>',
				'</li>',
				'</ul>',
				'<hr>',
				
			'</div>'
		].join('\n'));
	
	var display_section = $([
				'<div class="col-sm-9">',
				'<a href="#"><strong><i class="glyphicon glyphicon-dashboard"></i> My Dashboard</strong></a>',
				'<hr>',
				
				'<div class="jumbotron">',
					'<h1>Title</h1>',
					'<p>Some description.........</p>',
				'</div>',
				'<hr>',
				'<div class= "panel panel-default">',
					'<div class="panel-heading">',
						'<h4>Report</h4>',
					'</div>',
					'<div class="panel-body">',
					//	'<small>Decisions completed</small>',
						
					//	'<div class="progress">',
                    //            '<div class="progress-bar progress-bar-success" role="progressbar" aria-valuenow="72" aria-valuemin="0" aria-valuemax="100" style="width: 72%">',
                    //                '72 decisions Complete',
                    //            '</div>',
                    //    '</div>',
						'<small>Decisions in progress</small>',
						'<div class="progress">',
                                '<div class="progress-bar progress-bar-success" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="100" style="width: 5%">',
                                    '4 decisions in progress',
                                '</div>',
                            '</div>',
					'</div>',
				'</div>',
				
				'<div class="list-group">',
					
						
					'<a href="#" class="list-group-item active">',
							'Decisions In Progress',
					'</a>'
	].join('\n'));
	
	for(var i = 0; i < decisions_inProgress.length; i++){
		display_section.append('<a href="#" class="list-group-item">' + decisions_inProgress[i] + '</a>')
	}
					
	display_section.append('</div>' + '</div> \n')
	
	div_row.append(nav_section)
	div_row.append(display_section)
	div_dashboard.append(div_row)
	var footer = $([
	    '<div >',
        '<hr />',
        '<footer id="portfolio">',
            '<p>DECISION TOOL</p>',
        '</footer>',
    '</div>',
    '</div>',
	].join('\n'))
	footer.appendTo('body')
	
	
	
	$("a").attr("aria-expanded","true");
	$("a").click(function(){
		$(this).find('i#arrow_change').toggleClass('glyphicon-chevron-right glyphicon-chevron-down');
	});
	/*
	$("a").click(function(){
		$('#my-own-style').load("/static/person.html");
	}
	*/
} 





