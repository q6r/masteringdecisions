function main(body)
{
    get_text("/logout", function (result) { });
	window.location.replace("/login.html");
}