$(document).ready(function(){
    
    $("#spinner").hide();

    $("form").submit(function( event ){

        event.preventDefault();

        $("#spinner").show();

        $.ajax({
          type: "POST",
          url: "http://192.168.1.91/contact",
          data: $("#contact_form").serialize(),
          dataType: 'json',
          success: function(data, textStatus) {
                  $("#spinner").hide();

                  var is_error = data.Error;

                  var message_text=JSON.stringify(data.Message);
                  var message_title='<span class="message-header">Confirmation</span>';

                  message_text='<span class="message-success-text">' + message_text + '</span>';

                  if (is_error) {
                    message_title='<span class="message-header">Error</span>';
                    message_text='<span class="message-error-text">' + message_text + '</span>';
                  }

                  jquery_alert(message_title, message_text);
              },
          error: function(xhr, textStatus, errorThrown) {
                  $("#spinner").hide();
                          
                jquery_alert('<span class="">Error</span>',"<span class=\"\"> Form was not successfully submitted. Can you please try again.\n\n" + errorThrown + '</span>')
              }

        });

        document.getElementById("contact_form").reset();

    });

    function jquery_alert(title, content) {
        $.alert({
            title: title,
            content: content,
        });
    }

});
