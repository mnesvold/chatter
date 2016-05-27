(function(GT) {
  'use strict';

  function ChatClient(onrecv) {
    var url = 'ws://' + document.location.host + '/chat'
    this._ws = new WebSocket(url, []);
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
    var payload = JSON.stringify(msg);
    this._ws.send(payload);
  };

  document.addEventListener('DOMContentLoaded', function() {
    var $form = document.getElementById('message-form');
    var $input = document.getElementById('message-input');
    var $nickname = document.getElementById('nickname');
    var $transcript = document.getElementById('transcript');

    var client = new ChatClient(receiveMessage);

    function appendMessage(nick, msg) {
      var $p = document.createElement('p');
      var $nick = document.createElement('span');
      $nick.innerText = nick;
      $nick.classList.add('nickname');
      $p.appendChild($nick);
      $p.appendChild(document.createTextNode(msg));
      $transcript.appendChild($p);
    }

    function receiveMessage(jsonPayload) {
      var payload = JSON.parse(jsonPayload);
      var nick = payload.nickname || "Anonymous";
      var msg = payload.message;
      appendMessage(nick, msg);
    }

    $form.addEventListener('submit', function(ev) {
      var msg = $input.value;
      var nick = $nickname.value;
      $input.value = "";
      client.send({ nickname: nick, message: msg });
      ev.preventDefault();
    });

    $input.focus();
  });

})(window.GoTalk);
