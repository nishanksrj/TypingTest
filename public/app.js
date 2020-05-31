
// get the references of the page elements.
var formRoll = document.getElementById('form-roll');
var form = document.getElementById('form-msg');
var txtroll = document.getElementById('roll');
var txtMsg = document.getElementById('msg');
var listMsgs = document.getElementById('messages');
var socketStatus = document.getElementById('status');

var wordsElem = document.getElementById("wordsDisplay");
var charsElem = document.getElementById("charsDisplay");
var speedElem = document.getElementById("speedDisplay");

// variable to capture the starting time (Initialize it when the user starts typing)
var start;

// Initiate a web socket connection
var socket = new WebSocket('ws://' + window.location.host + '/ws');

// clear the input and text fields on page refresh
window.onload = function(e){
	
	txtMsg.value = '';
	txtroll.value = '';
}

// make the Roll Number input disappear and display the text field
formRoll.onsubmit = function(e){
	e.preventDefault();
	formRoll.style.display = "none";
	document.getElementById("rollDisplay").innerHTML = txtroll.value;
	form.style.display = "block";
	txtMsg.focus();
}


// Display the change in socket status
socket.onopen = function(event) {
	socketStatus.innerHTML = 'Connected to: ' + event.currentTarget.URL;
	socketStatus.className = 'open';
};

// Display the change in socket status
socket.onerror = function(error) {
	socketStatus.innerHTML = 'Error:  ' + error ;
	socketStatus.className = 'close';
};


// Display the change in socket status
socket.onclose = function(){
	socketStatus.innerHTML = 'Connection closed.';
};


// Send the data to the server through websocket when input changes
txtMsg.oninput = function() {

	// Initialize the time variable
	if(start==undefined){
		start = performance.now();
	}

	// capture the current time
	var curr = performance.now();
	
	// diff between curr and start in seconds
	var diff = (curr - start)/1000.0;

	// get the roll number from roll number field
	var roll = txtroll.value;

	// Recovering the message of the textarea.
    	var str = txtMsg.value;
	str = str.replace(/(^\s*)|(\s*$)/gi,"");
	str = str.replace(/[ ]{2,}/gi," ");
        str = str.replace(/\n /,"\n");
	var chars = str.length;
	var words = str.split(' ').length;
	
	var speed = parseInt(words*60/diff);

	// update the stats
	wordsElem.innerHTML = words;
	charsElem.innerHTML = chars;
	speedElem.innerHTML = speed;

    	// Sending the msg via WebSocket.
    	socket.send(JSON.stringify({roll: roll, message: str, speed: speed, words: words, chars: chars}));

};

