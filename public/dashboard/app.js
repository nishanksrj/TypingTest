
// get the references of the page elements.
var form = document.getElementById('form-msg');
var stats = document.getElementById('stats');
var socketStatus = document.getElementById('status');

// Creating a new WebSocket connection.
 var socket = new WebSocket('ws://' + window.location.host + '/dash');


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
	


// Render the recieved data 
socket.onmessage = function(event) {

	// parse and extract the recieved data
    	var data = JSON.parse(event.data);
	roll = data['roll'];
	speed = data['speed'];
	words = data['words'];
	chars = data['chars'];

	// get the element to be updated
	var elem = document.getElementById(roll);
	if(elem){
		// if element exist, render the data
		elem.innerHTML = '<th>' + roll + ': </th><th>' + words + '</th><th>' + chars +'</th><th>' + speed + '</th>'; 
	}
	else{
		// otherwise create the element and display the data
		stats.innerHTML+= '<tr class="received" id="' + data['roll'] + '"><th>' + roll + ': </th><th>' + words + '</th><th>' + chars +'</th><th>' + speed + '</th></tr>';
	}
	return false;
};

