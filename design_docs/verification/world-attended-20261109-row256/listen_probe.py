import socket, time, subprocess, sys
s = socket.socket(); s.bind(("127.0.0.1", 0)); s.listen(16); s.settimeout(1.0)
port = s.getsockname()[1]; print("port", port, flush=True)
end = time.time() + float(sys.argv[1])
while time.time() < end:
    try:
        c, a = s.accept()
    except socket.timeout:
        continue
    c.settimeout(0.5)
    try: b = c.recv(200)
    except Exception: b = b""
    who = subprocess.run(["lsof", "-nP", "-iTCP@127.0.0.1:%d" % a[1], "-sTCP:ESTABLISHED"], capture_output=True, text=True).stdout
    print(time.strftime("%H:%M:%S"), a, b[:120], "|", " ".join(l.split()[0]+"/"+l.split()[1] for l in who.splitlines()[1:]), flush=True)
    c.close()
