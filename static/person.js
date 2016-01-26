function main(body)
{
    $('title')
        .html('person information');

    $('<p>')
        .html('Person is : ' + JSON.stringify(body))
        .appendTo('body');
}
