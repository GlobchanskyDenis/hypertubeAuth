var serverIP = "localhost"
var serverPort = "4000"
var token = ""

function AuthBasic() {
	var email = document.forms['authBasic']['email'].value
	var passwd = document.forms['authBasic']['passwd'].value
	var authRaw = btoa(encodeURI(email)+":"+encodeURI(passwd))
	console.log("tx: " + authRaw)

	let xhr = new XMLHttpRequest();
	xhr.open("GET", "http://"+serverIP+":"+serverPort+"/api/auth/basic");
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
			document.getElementById("errorField").innerHTML += "access token is empty. Its not a valid case";
			return
		}
		document.getElementById("responseField").innerHTML = "access_token=hidden"
		console.log("rx token: " + requestAsync["access_token"]);
		console.log("rx user profile: " + requestAsync["profile"]);
		document.forms['authBasic']['email'].value = "";
		document.forms['authBasic']['passwd'].value = "";
		document.token = requestAsync["access_token"]
	}
	xhr.onerror = function () {
		console.log("onError event")
	}
}

function ProfileCreate() {
	var email = document.forms['profileCreate']['email'].value
	var pass = document.forms['profileCreate']['passwd'].value
	var username = document.forms['profileCreate']['username'].value
	var user = {
		email: email,
		passwd: pass,
		username: username
	}
	var request = JSON.stringify(user)
	let xhr = new XMLHttpRequest();
	xhr.open("PUT", "http://"+serverIP+":"+serverPort+"/api/profile/create");
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
		document.forms['profileCreate']['email'].value = "";
		document.forms['profileCreate']['passwd'].value = "";
		document.forms['profileCreate']['username'].value = "";
	}
	xhr.onerror = function () {
		console.log("onError event")
	}
}

function ProfilePatch() {
	var fname = document.forms['profilePatch']['firstName'].value
	var lname = document.forms['profilePatch']['lastName'].value
	var username = document.forms['profilePatch']['username'].value
	var imageBody = document.getElementById('avatar').src;
	var user = {};
	if (fname != "") {
		user.firstName = fname;
	}
	if (lname != "") {
		user.lastName = lname;
	}
	if (username != "") {
		user.username = username;
	}
	if (imageBody != "") {
		user.imageBody = imageBody;
	}
	var request = JSON.stringify(user)
	let xhr = new XMLHttpRequest();
	xhr.open("PATCH", "http://"+serverIP+":"+serverPort+"/api/profile/patch");
	xhr.setRequestHeader("access_token", document.token)
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
		document.getElementById("responseField").innerHTML = "user profile was patched"
		document.forms['profilePatch']['firstName'].value = "";
		document.forms['profilePatch']['lastName'].value = "";
		document.forms['profilePatch']['username'].value = "";
	}
	xhr.onerror = function () {
		console.log("onError event")
	}
}

function readURL(input) {
	if (input.files && input.files[0]) {
		var reader = new FileReader();
		image = document.getElementById('avatar');
		image.style.display = "block";
		reader.onload = function (e) {
			image.setAttribute('src', e.target.result);
		};
		reader.readAsDataURL(input.files[0]);
	}
	image_statut = true;
}