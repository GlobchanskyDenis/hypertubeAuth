var serverIP = "localhost"
var serverPort = "4000"
var token = ""

function AuthUser() {
	var mail = document.forms['auth']['mail'].value
	var pass = document.forms['auth']['pass'].value
	var authRaw = btoa(encodeURI(mail)+":"+encodeURI(pass))
	console.log("tx: " + authRaw)

	let xhr = new XMLHttpRequest();
	xhr.open("GET", "http://"+serverIP+":"+serverPort+"/user/auth/basic");
	xhr.setRequestHeader("Authorization", "Basic "+authRaw);
	xhr.send()

	xhr.onload = function () {
		if (xhr.response) {
			var requestAsync = JSON.parse(xhr.response);
		} else {
			var requestAsync = "";
		}
		console.log("rx: " + xhr.status + " : " + xhr.response);
		if (xhr.status != 200) {
			document.getElementById("errorField").innerHTML = "Что-то пошло не так: " + xhr.status + " : "
			document.getElementById("errorField").innerHTML += ((requestAsync.error) ? requestAsync.error : xhr.statusText)
			document.getElementById("responseField").innerHTML = "";
			return
		}
		if (!requestAsync) {
			document.getElementById("errorField").innerHTML = "Empty request. Its not a valid case";
			document.getElementById("responseField").innerHTML = "";
			return
		}
		document.getElementById("errorField").innerHTML = ""
		if (!requestAsync.access_token) {
			document.getElementById("errorField").innerHTML += "uid or token are empty. Its not a valid case";
			return
		}
		document.getElementById("responseField").innerHTML = "access_token=hidden"
		console.log("rx token: " + requestAsync["access_token"]);
		console.log("rx user profile: " + requestAsync["profile"]);
		document.forms['auth']['mail'].value = "";
		document.forms['auth']['pass'].value = "";
		document.token = requestAsync["access_token"]
	}
	xhr.onerror = function () {
		console.log("onError event")
	}
}

function oAuth42() {

	let xhr = new XMLHttpRequest();
	xhr.open("GET", "https://api.intra.42.fr/oauth/authorize?"+
	"client_id=96975efecfd0e5efee67c9ac4cc350ac9372ae559b2fb8a08feba6841a33fb53"+
	"&redirect_uri=http://localhost:4000/user/auth/oauth42"+
	"&scope=public"+
	"&state=bdcbe28874ab05962b50430b1466a8ebcbda45ba8c3c1beee600699478ad2a4d"+
	"&response_type=code", true);
	// xhr.setRequestHeader("Access-Control-Allow-Origin", "*");
	xhr.send()

	xhr.onload = function () {
		if (xhr.response) {
			var requestAsync = JSON.parse(xhr.response);
		} else {
			var requestAsync = "";
		}
		console.log("rx: " + xhr.status + " : " + xhr.response);
		if (xhr.status != 200) {
			document.getElementById("errorField").innerHTML = "Что-то пошло не так: " + xhr.status + " : "
			document.getElementById("errorField").innerHTML += ((requestAsync.error) ? requestAsync.error : xhr.statusText)
			document.getElementById("responseField").innerHTML = "";
			return
		}
		if (!requestAsync) {
			document.getElementById("errorField").innerHTML = "Empty request. Its not a valid case";
			document.getElementById("responseField").innerHTML = "";
			return
		}
		document.getElementById("errorField").innerHTML = ""
		if (!requestAsync.access_token) {
			document.getElementById("errorField").innerHTML += "access_token are empty. Its not a valid case";
			return
		}
		console.log("rx token: " + requestAsync["access_token"]);
		console.log("rx user profile: " + requestAsync["profile"]);
		document.getElementById("responseField").innerHTML = "access_token=hidden"
		document.token = requestAsync["access_token"]
	}
	xhr.onerror = function () {
		console.log("onError event")
	}
}

function RegUser() {
	var mail = document.forms['reg']['mail'].value
	var pass = document.forms['reg']['pass'].value
	var fname = document.forms['reg']['first_name'].value
	var lname = document.forms['reg']['last_name'].value
	var displayname = document.forms['reg']['displayname'].value
	var request = JSON.stringify({ "email": mail, "passwd": pass })
	var user = {
		email: mail,
		passwd: pass,
		first_name: fname,
		last_name: lname,
		displayname: displayname
	}
	var request = JSON.stringify(user)
	let xhr = new XMLHttpRequest();
	xhr.open("PUT", "http://"+serverIP+":"+serverPort+"/user/create/basic");
	console.log("tx: " + request)
	xhr.send(request);
	xhr.onload = function () {
		if (xhr.response) {
			var requestAsync = JSON.parse(xhr.response);
		} else {
			var requestAsync = "";
		}
		console.log("rx: " + xhr.status + " : " + xhr.response);
		if (xhr.status != 200) {
			document.getElementById("errorField").innerHTML = "Что-то пошло не так: " + xhr.status + " : "
			document.getElementById("errorField").innerHTML += ((requestAsync.error) ? requestAsync.error : xhr.statusText)
			document.getElementById("responseField").innerHTML = "";
			return
		}
		document.getElementById("errorField").innerHTML = ""
		document.getElementById("responseField").innerHTML = "registration was done. Check your email"
		document.forms['reg']['mail'].value = "";
		document.forms['reg']['pass'].value = "";
	}
	xhr.onerror = function () {
		console.log("onError event")
	}
}