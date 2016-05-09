(function(GT) {
  'use strict';

  function ChatClient(onrecv) {
    var url = 'ws://' + document.location.host + '/chat'
    this._ws = new WebSocket(url, 'go-talk/1.0');
    this._ws.addEventListener('close', function() {
      console.warn('WS connection closed');
    });
    this._ws.addEventListener('error', function(err) {
      console.error(err);
    });
    this._ws.addEventListener('message', function(evt) {
      onrecv(evt.data);
    });
  }

  ChatClient.prototype.send = function(msg) {
    var payload = JSON.stringify({ 'message': msg });
    this._ws.send(payload);
  };

  document.addEventListener('DOMContentLoaded', function() {
    var $form = document.getElementById('message-form');
    var $input = document.getElementById('message-input');
    var $transcript = document.getElementById('transcript');

    var client = new ChatClient(receiveMessage);

    function appendMessage(msg) {
      var $p = document.createElement('p');
      $p.innerText = msg;
      $transcript.appendChild($p);
    }

    function receiveMessage(jsonPayload) {
      var payload = JSON.parse(jsonPayload);
      var msg = payload.message;
      appendMessage(msg);
    }

    $form.addEventListener('submit', function() {
      var msg = $input.value;
      $input.value = "";
      client.send(msg);
    });

    $input.focus();
  });

})(window.GoTalk);
