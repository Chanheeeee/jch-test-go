from flask import Flask, request
import ssl


app = Flask(__name__)

@app.route('/', methods=['POST'])
def webhook():    
    # validate from INPUT data
    result = validate(request.json)
 
    return {
      "kind": "AdmissionReview",
      "apiVersion": "admission.k8s.io/v1beta1",
      "response": {
        "allowed": result,
        "status": {
          "reason": "Pod create not allowed"
        }
      }
}


def validate(review):
    # denying all Pod creating
    if (review['request']['object']['kind'] == 'Pod') and \
        (review['request']['operation'] == 'CREATE'):
        return False  # Deny
    return True       # Accept


##################################
# Webhook needs to serve TLS
##################################
context = ssl.SSLContext(ssl.PROTOCOL_TLS)
context.load_verify_locations('./ca.crt')
context.load_cert_chain('./server.crt', './server.key')

app.run(host='0.0.0.0', debug=True, ssl_context=context)