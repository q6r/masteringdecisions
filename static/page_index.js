function main(body)
{
	//Add CSS
	$.loadCSS('static/css/index.css');
	
	buildTemplate();
	
	buildHome();
	
} 

function buildTemplate() {
	var decisions_inProgress = ["decision1", "decision2", "decision3","decision4"]
	var decisions_completed = ["decision5", "decision6", "decision7","decision8"]
	
	//nav section
	var nav = $('<nav>').addClass('navbar navbar-inverse navbar-fixed-top').appendTo('body')
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
	div_nav_header.append($('<a class = "navbar-brand"><img id="logo" src="../static/images/logo.png">'))
	
	var div_collapse = $('<div class="collapse navbar-collapse" id="myNavbar">')
	
	var nav_ul1 = $('<ul class="nav navbar-nav"><li> <a onclick="buildHome()">Dashboard</a></li></ul>')
	var nav_ul2 = $([
	'<ul class="nav navbar-nav navbar-right">',
		'<li class="dropdown">',
			'<a class="dropdown-toggle" role="button" data-toggle="dropdown" href="#" aria-expanded="false"><i class="glyphicon glyphicon-user"></i><span id="userName">' + '</span><span class="caret"></span></a>',
				'<ul id="g-account-menu" class="dropdown-menu" role="menu">',
					'<li><a onclick="buildEditUser()">Edit Profile</a></li>',
				'</ul>',
		'</li>',
		'<li><a href="/logout.html"><i class="glyphicon glyphicon-lock"></i> Logout</a></li>',
	'</ul>'
	].join("\n"));
				
				
	div_collapse.append(nav_ul1)
	div_collapse.append(nav_ul2)
	div_container.append(div_collapse)
	
	//dashboard section
	
	var div_dashboard = $('<div class="container-fluid" style="margin-top:75px;">').appendTo('body')
	
	var div_row = $('<div class="row-fluid">')
	
	var nav_section = $([
			
			'<div class="col-sm-3" >',
				'<a>',
					'<strong><i class="glyphicon glyphicon-wrench"></i> Tools</strong>',
				'</a>',
				'<hr>',
				'<ul class="nav nav-stacked">',
				'<li class="nav-header"> <a  data-toggle="collapse" data-target="#userMenu" aria-expanded="false" class="collapsed">Decisions <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
                    '<ul class="nav nav-stacked collapse" id="userMenu" aria-expanded="false" style="height: 0px;">',
						
                        '<li class="active"> <a ><i class="glyphicon glyphicon-asterisk"></i> New Decision</a></li>',
						'<li class="nav-header"> <a  data-toggle="collapse" data-target="#userMenu3" aria-expanded="false" class="collapsed">In Progress <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
							'<ul class="nav nav-stacked collapse" id="userMenu3" aria-expanded="false" style="height: 0px;">',
								'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Decision1 </a></li>',
								'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Decision2 </a></li>',
								'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Decision3 </a></li>',
								'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Decision4 </a></li>',
							'</ul>',
						'</li>',
                       '<li class="nav-header"> <a  data-toggle="collapse" data-target="#userMenu4" aria-expanded="false" class="collapsed">Completed <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
							'<ul class="nav nav-stacked collapse" id="userMenu4" aria-expanded="false" style="height: 0px;">',
								'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Decision5 </a></li>',
								'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Decision6 </a></li>',
								'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Decision7 </a></li>',
								'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Decision8 </a></li>',
							'</ul>',
						'</li>',
                    '</ul>',
                '</li>',
				'<li class="nav-header"> <a  data-toggle="collapse" data-target="#Menu2" aria-expanded="false" class="collapsed">Settings <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
					
                    '<ul class="nav nav-stacked collapse" id="Menu2" aria-expanded="fasle">',
                        '<li class="active"> <a ><i class="glyphicon glyphicon-home"></i> Home</a></li>',
                        
                        //'<li><a ><i class="glyphicon glyphicon-cog"></i> Options</a></li>',
                        
                        '<li><a href="/static/person.html"><i class="glyphicon glyphicon-user"></i> New User</a></li>',
                        '<li><a ><i class="glyphicon glyphicon-list-alt"></i> Report</a></li>',
                        
                        '<li><a ><i class="glyphicon glyphicon-off"></i> Logout</a></li>',
						
                    '</ul>',
				'</li>',
				'</ul>',
				'<hr>',
				
			'</div>'
		].join('\n'));
	
	var display_section = $('<div class="col-sm-9" id="content">');
	
	div_row.append(nav_section)
	div_row.append(display_section)
	div_dashboard.append(div_row)

	updateUserText();
	
	$("a").attr("aria-expanded","true");
	$("a").click(function(){
		$(this).find('i#arrow_change').toggleClass('glyphicon-chevron-right glyphicon-chevron-down');
	});
}

function clearContent() {
	$('#content').empty();
}

function updateUserText() {
	get_text("/whoami", function (result) {
		get_text("/person/"+result['person_id']+"/info", function (result) {
			$('#userName').text(' ' + result['person']['name_first'] + ' ' + result['person']['name_last']);
		});
	});
}

function buildHome() {
	$('title').html('Decision Home');
	
	clearContent();
	
	var decisions_inProgress = ["decision1", "decision2", "decision3","decision4"]
	var decisions_completed = ["decision5", "decision6", "decision7","decision8"]
	
	var display_section = $([
				'<div>',
				'<strong><i class="glyphicon glyphicon-dashboard"></i> My Dashboard</strong>',
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
					
						
					'<a  class="list-group-item active">',
							'Decisions In Progress',
					'</a>'
	].join('\n'));
	
	for(var i = 0; i < decisions_inProgress.length; i++){
		display_section.append('<a  class="list-group-item">' + decisions_inProgress[i] + '</a>')
	}
					
	display_section.append('</div>' + '</div> \n')
	
	display_section.appendTo('#content');
}

function buildEditUser() {
	$('title').html('Update User!');

	clearContent();
	
	$('<strong><i class="glyphicon glyphicon-cog"></i> Edit Person</strong><hr/>').appendTo('#content');
	
	var wrapper = $('<div>').css('max-width','500px').appendTo('#content');
	var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return updateUser()').appendTo(wrapper);
		$('<h3>').addClass('form-signin-heading').text('Update Details').appendTo(form);
		
		$('<div>').attr('id','success').addClass('alert alert-success').appendTo(form);
		$('#success').hide();
		
		$('<input type="text" />').addClass('form-control')
			.attr('name', 'username')
			.attr('placeholder', 'Email Address')
			.attr('id', 'username')
			.appendTo(form);
		
		$('<input type="text" />').addClass('form-control')
			.attr('name', 'firstname')
			.attr('placeholder', 'First Name')
			.attr('id', 'firstname')
			.appendTo(form);
		
		$('<input type="text" />').addClass('form-control')
			.attr('name', 'lastname')
			.attr('placeholder', 'Last Name')
			.attr('id', 'lastname')
			.appendTo(form);
		
		$('<h3>').addClass('form-signin-heading').text('Password').appendTo(form);
		
		$('<div>').attr('id','error').addClass('alert alert-danger').appendTo(form);
		$('#error').hide();
		
		$('<input type="password" />').addClass('form-control')
			.attr('name', 'password')
			.attr('placeholder', 'Password')
			.attr('id', 'password')
			.appendTo(form);
			
		$('<input type="password" />').addClass('form-control')
			.attr('name', 'password2')
			.attr('placeholder', 'Password')
			.attr('id', 'password2')
			.appendTo(form);

		$('<hr/>').appendTo(form);
		$('<button>').addClass('btn btn-lg btn-primary btn-block').attr('type', 'submit').text('Submit').appendTo(form);
		
		get_text("/whoami", function (result) {
			get_text("/person/"+result['person_id']+"/info", function (result) {
				$('#username').val(result['person']['email']);
				$('#firstname').val(result['person']['name_first']);
				$('#lastname').val(result['person']['name_last']);
			});
		});
}

function updateUser() {
	$('#error').hide();
	
	if(document.getElementById('password').value != document.getElementById('password2').value) {
		$('#error').html('<b>Error:</b> Passwords do not match!');
		$('#error').show();
	}
	else {
		if(document.getElementById('password').value == "") {
			new_info = {
				"email":document.getElementById('username').value,
				"name_first":document.getElementById('firstname').value,
				"name_last":document.getElementById('lastname').value
				}
		}
		else {
			new_info = {
					"email":document.getElementById('username').value,
					"pw_hash":document.getElementById('password').value,
					"name_first":document.getElementById('firstname').value,
					"name_last":document.getElementById('lastname').value
					}
		}
		
		get_text("/whoami", function (result) {
			put_text("/person/"+result['person_id'], JSON.stringify(new_info), function (result) {
				//alert(JSON.stringify(result));
				document.getElementById('password').value = "";
				document.getElementById('password2').value = "";
				$('#success').html('<b>Update successful!</b>');
				$('#success').show();
				updateUserText();
			});
		});
	}
	
	return false;
}
