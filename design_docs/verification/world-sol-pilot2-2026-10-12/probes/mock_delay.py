"""Fake Responses-API endpoint on 127.0.0.1: records request bodies, answers one short message. No provider call."""
import http.server,json,sys,os,itertools
OUT=sys.argv[1]; PORT=int(sys.argv[2]); n=itertools.count()
def sse(ev): return f"event: {ev['type']}\ndata: {json.dumps(ev)}\n\n".encode()
class H(http.server.BaseHTTPRequestHandler):
    def log_message(self,*a): pass
    def do_GET(self):
        self.send_response(200); self.send_header('Content-Type','application/json'); self.end_headers()
        self.wfile.write(b'{"object":"list","data":[],"models":[]}')
    def do_POST(self):
        body=self.rfile.read(int(self.headers.get('Content-Length',0)))
        i=next(n); open(os.path.join(OUT,f'req-{i}{self.path.replace("/","_")}.json'),'wb').write(body)
        import time as _t; _t.sleep(float(os.environ.get('MOCK_DELAY','0')))
        self.send_response(200); self.send_header('Content-Type','text/event-stream'); self.end_headers()
        JS=os.environ.get('MOCK_EXEC')
        if i==0 and JS:
            msg={'type':'custom_tool_call','id':'ctc_1','call_id':'call_1','name':'exec','input':JS,'status':'completed'}
        else:
          msg={'type':'message','id':'msg_1','role':'assistant','status':'completed','content':[{'type':'output_text','text':'done','annotations':[]}]}
        resp={'id':'resp_1','object':'response','status':'completed','output':[msg],'usage':{'input_tokens':1,'output_tokens':1,'total_tokens':2,'input_tokens_details':{'cached_tokens':0},'output_tokens_details':{'reasoning_tokens':0}}}
        for ev in ({'type':'response.created','response':{**resp,'status':'in_progress','output':[]}},
                   {'type':'response.output_item.done','output_index':0,'item':msg},
                   {'type':'response.completed','response':resp}):
            self.wfile.write(sse(ev))
http.server.ThreadingHTTPServer(('127.0.0.1',PORT),H).serve_forever()
