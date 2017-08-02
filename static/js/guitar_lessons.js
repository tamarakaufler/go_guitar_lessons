$(document).ready(function(){
    
    $("#gl_spinner").hide();

    $("form").submit(function( event ){

        event.preventDefault();

        $("#gl_spinner").show();

        $.ajax({
          type: "POST",
          url: "/contact",
          data: $("#contact_form").serialize(),
          dataType: 'json',
          success: function(data, textStatus) {
                  $("#gl_spinner").hide();

                  var is_error = data.Error;

                  var message=JSON.stringify(data.Message);
                  var message_title='<span class="message-header">Confirmation</span>';

                  message_text='<span class="message-success-text">' + message + '</span>';

                  if (is_error) {
                      if (message.match("Redirect")) {

                          grecaptcha.reset();
                          document.getElementById("contact_form").reset();

                          window.location('/');
                          return;

                      } else {
                          message_title='<span class="message-header">Error</span>';
                          message_text='<span class="message-error-text">' + message + '</span>';
                      }
                  }

                  jquery_alert(message_title, message_text);

                  grecaptcha.reset();
                  document.getElementById("contact_form").reset();
              },
          error: function(xhr, textStatus, errorThrown) {
                  $("#gl_spinner").hide();

                  jquery_alert('<span class="">Error</span>','<span class="">Form was not successfully submitted. Can you please try again.</span>')

                  grecaptcha.reset();
                  document.getElementById("contact_form").reset();
              }

        });

    });

    function jquery_alert(title, content) {
        $.alert({
            title: title,
            content: content
        });
    }
});

function initMap() {

    var housePos = {lat: 51.777, lng: -0.397};
    var map = new google.maps.Map(document.getElementById("map"), {
      zoom: 13,
      center: housePos
    });

    var marker = new google.maps.Marker({
      position: housePos,
      map: map
    });

}

