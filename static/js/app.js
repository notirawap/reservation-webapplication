function Prompt() {
        let success = function(c) {
          const {
            msg = "",
            title = "",
            footer = "",
          } = c;

          Swal.fire({
            icon: "success",
            title: title,
            text: msg,
            footer: footer,
          });
        }

        let error = function(c) {
          const {
            msg = "",
            title = "",
            footer = "",
          } = c;

          Swal.fire({
            icon: "error",
            title: title,
            text: msg,
            footer: footer,
          });
        }

        // custom popup dialog
        async function custom(c) {
          const {
            icon = "",
            msg = "",
            showConfirmButton = true,
            showCancelButton = true,
          } = c;
          
          const result = await Swal.fire({
            icon: icon,
            html: msg,
            showCancelButton: showCancelButton,
            showConfirmButton: showConfirmButton,
          });

          // Process code after clicking the submit button
          // Generalize function
          if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
              if (result.value !== "") {
                if (c.callback !== undefined) {
                  c.callback(result); // execute by passing the result back
                }
              } else {
                c.callback(false); // execute if the result is ""
              }
            } else {
              c.callback(false); // execute if the result is dismissed (cancel button clicked)
            }
          };
        }

        // When prompt get a request for toast, it returns toast.
        return {
          toast: toast,
          success: success,
          error: error,
          custom: custom,
        }
}