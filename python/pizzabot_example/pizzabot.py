#!/usr/bin/env python
from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
import SocketServer
import re
import requests
import json
import multiprocessing
import string
import random
import urllib

# Integration data from developers console
CLIENT_ID = '<CLIENT_ID>'
CLIENT_SECRET = '<CLIENT_SECRET>'
REDIRECT_URI = '<REDIRECT_URI>'

# BOT Agent data
BOT_AGENT_WEBHOOKS_URL = '<BOT_AGENT_WEBHOOKS_URL>'
BOT_AGENT_WEBHOOKS_SECRET = ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(10))
BOT_AGENT_NAME = 'Pizza BOT'

# API urls
ACCOUNTS_TOKEN_URL = 'https://accounts.chat.io/token'
ACCOUNTS_INFO_URL = 'https://accounts.chat.io/info'
CONFIGURATION_API_URL = 'https://api.chat.io/configuration/v0.4'
AGENT_API_URL = 'https://api.chat.io/agent/v0.4/action'


class Bot:
    id = ''
    license = 0
    access_token = ''
    initialized = False
    
    # Helpers
    def sendConfigurationAPIRequest(self, path, payload):
        result = requests.post('{0}{1}'.format(CONFIGURATION_API_URL, path),
            json=payload,
            headers={'Authorization': 'Bearer {0}'.format(self.access_token)})
            
        if result.status_code != 200:
            print "Configuration API request to {0} error: code: {1}, text: {2}".format(path, result.status_code, result.text)
            return ''
        
        return result.json()

    def sendAgentAPIRequest(self, path, payload):
        result = requests.post('{0}{1}'.format(AGENT_API_URL, path),
            json={
                'author_id': self.id,
                'payload': payload
            },
            headers={'Authorization': 'Bearer {0}'.format(self.access_token)})
            
        if result.status_code != 200:
            print "Agent API request to {0} error: code: {1}, text: {2}".format(path, result.status_code, result.text)
            return ''
        
        return result.json()
        
    
    # BOT agent initialization
    def createBotAgent(self):
        print "Creating BOT Agent for license {0}".format(self.license)
        response_data = self.sendConfigurationAPIRequest('/agents/create_bot_agent', {
                        'name': BOT_AGENT_NAME,
                        'status': 'not accepting chats',
                        'webhooks': {
                            'url': BOT_AGENT_WEBHOOKS_URL,
                            'secret_key': BOT_AGENT_WEBHOOKS_SECRET,
                            'actions': [{
                                'name': 'incoming_event'
                            },{
                                'name': 'incoming_chat_thread'
                            }]
                        }
                    })

        if response_data == '':
            print "BOT Agent creation failed.".format(self.id)
            return ''
        
        bot_agent_id = response_data['bot_agent_id']
        
        print "BOT Agent created with ID {0}, integration enabled.".format(bot_agent_id)
        
        return bot_agent_id
    
    def __init__(self, license, access_token):
        self.license = license
        self.access_token = access_token

        self.id = self.createBotAgent()
        
        if self.id != '':
            self.initialized = True
            

    # BOT actions
    def joinChat(self, chat_id):
        print 'Joining chat: {0}'.format(chat_id)
        response_data = self.sendAgentAPIRequest('/join_chat', {
                    'chat_id': chat_id,
                    'agent_ids': [self.id]
                })
            
        if response_data == '':
            print "Can not join chat: {0}".format(result.text)
            return False
        
        return True

    def sendMessageToChat(self, chat_id, message):
        print 'Sending message to chat: {0}'.format(chat_id)
        response_data = self.sendAgentAPIRequest('/send_event', {
                    'chat_id': chat_id,
                    'event': {
                        'type': 'message',
                        'text': message,
                        'recipients': 'agents'
                    }
                })

        if response_data == '':
            print "Can send message to chat {0}".format(result.text)
            return False
            
        return True
             
    def leaveChat(self, chat_id):
        print 'Leaving chat: {0}'.format(chat_id)
        response_data = self.sendAgentAPIRequest('/remove_from_chat', {
                    'chat_id': chat_id,
                    'agent_ids': [self.id]
                })

        if response_data == '':
            print "Can not leave chat {0}".format(result.text)
            return False
        
        return True

    def onIncomingKeyWord(self, chat_id):
        if self.joinChat(chat_id):
            self.sendMessageToChat(chat_id, 'Pizza is on the way!')
            self.leaveChat(chat_id)

    # Webhooks handling
    def manageIncomingEvent(self, payload):
        if 'message' in payload['event'] and 'pizza' in payload['event']['message']:
            self.onIncomingKeyWord(payload['chat_id'])

    def manageIncomingChatThread(self, payload):
        for event in payload['chat']['thread']['events']:
            if 'message' in event and 'pizza' in event['message']:
                self.onIncomingKeyWord(payload['chat']['id'])

    def incomingWebhook(self, payload):
        if payload['secret_key'] != BOT_AGENT_WEBHOOKS_SECRET:
            print "Wrong incoming secret key: {0}".format(payload['secret_key'])
            return
            
        if payload['action'] == 'incoming_event':
            self.manageIncomingEvent(payload['data'])
        elif payload['action'] == 'incoming_chat_thread':
            self.manageIncomingChatThread(payload['data'])

class S(BaseHTTPRequestHandler):    
    manager = multiprocessing.Manager()
    SHARED = manager.dict()

    def sendCannotStartBotIntegration(self):
        self.setHeadersOK()
        self.wfile.write("<html><body>")
        self.wfile.write("Can not enable BOT integration")
        self.wfile.write("</body></html>")        

    def sendAlreadyEnabledBotIntegration(self):
        self.setHeadersOK()
        self.wfile.write("<html><body>")
        self.wfile.write("BOT integration already enabled")
        self.wfile.write("</body></html>")        

    def sendSuccessfullyEnabledBotIntegration(self, license_id):
        self.setHeadersOK()
        self.wfile.write("<html><body>")
        self.wfile.write("BOT integration successfully enabled on license {0}!".format(license_id))
        self.wfile.write("</body></html>")        

    def sendOK(self):
        self.setHeadersOK()
        
    def exchangeCodeForAccessToken(self, code):
        print "Exchanging SSO code for access token"
        result = requests.post(ACCOUNTS_TOKEN_URL, data={'grant_type': 'authorization_code', 'code': code, 'client_id': CLIENT_ID, 'client_secret': CLIENT_SECRET, 'redirect_uri': REDIRECT_URI})
        if result.status_code != 200:
            self.sendCannotStartBotIntegration()
            return

        response_data = result.json()
        access_token = response_data['access_token']
        
        print "Getting access token info to get license number"
        
        result = requests.get(ACCOUNTS_INFO_URL, headers={'Authorization': 'Bearer {0}'.format(access_token)})
        if result.status_code != 200:
            self.sendCannotStartBotIntegration()
            return
        
        response_data = result.json()
        license_id = response_data['license_id']

        print "Creating BOT Agent object"
        
        self.SHARED['bot'] = Bot(license_id, access_token)
        if self.SHARED['bot'].initialized == False:
            del self.SHARED['bot']
            selfsendCannotStartBotIntegration()
            return
        
        self.sendSuccessfullyEnabledBotIntegration(license_id)
    
    def setHeadersOK(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

    def do_GET(self):
        if 'bot' in self.SHARED:
            self.sendAlreadyEnabledBotIntegration()
            return

        if self.path.startswith('/token'):
            result = re.search('code=([^&]+)', self.path)
            if result != None:
                self.exchangeCodeForAccessToken(urllib.unquote(result.group(1)).decode('utf8'))
                return
        
        self.sendOK()
            

    def do_HEAD(self):
        self.setHeadersOK()
        
    def do_POST(self):
        if 'bot' not in self.SHARED:
            self.sendOK()
            return
        
        if self.path == '/webhook':
            content_len = int(self.headers.getheader('content-length', 0))
            post_body = self.rfile.read(content_len)
            self.SHARED['bot'].incomingWebhook(json.loads(post_body))

        self.sendOK()
        
def run(server_class=HTTPServer, handler_class=S, port=80):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print 'Starting httpd...'
    httpd.serve_forever()

if __name__ == "__main__":
    from sys import argv

    if len(argv) == 2:
        run(port=int(argv[1]))
    else:
        run()
