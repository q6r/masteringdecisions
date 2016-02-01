function main(body)
{
    $('title')
        .html('title goes here');

    $('<p>')
        .html('Data : ' + JSON.stringify(body))
        .appendTo('body');
}
