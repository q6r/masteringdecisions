function main(body)
{
    $('title')
        .html('title goes here');

	$.each(body, function(index, jsonObject){
			$('<li>')
				.html(jsonObject["email"])
				.appendTo('body');
	});
}
