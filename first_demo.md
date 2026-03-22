### 3.3 Tạo script relay demo

Tạo file `/opt/rtcm-relay/relay_from_tcpdump.py` với nội dung sau:

```python
#!/usr/bin/env python3
import argparse
import socket
import struct
import subprocess
import sys
import time
from collections import defaultdict


def parse_pcap_stream(stream, target_port):
    gh = stream.read(24)
    if len(gh) < 24:
        return
    magic = struct.unpack('<I', gh[:4])[0]
    if magic not in (0xA1B2C3D4, 0xD4C3B2A1):
        raise RuntimeError('PCAP header không hợp lệ')

    while True:
        ph = stream.read(16)
        if len(ph) < 16:
            break
        ts_sec, ts_usec, incl_len, orig_len = struct.unpack('<IIII', ph)
        pkt = stream.read(incl_len)
        if len(pkt) < incl_len:
            break

        # Linux SLL2 link header: 20 bytes
        if len(pkt) < 20:
            continue
        proto = struct.unpack('>H', pkt[0:2])[0]
        if proto != 0x0800:
            continue
        ip = pkt[20:]
        if len(ip) < 20:
            continue

        ihl = (ip[0] & 0x0F) * 4
        if len(ip) < ihl + 20:
            continue
        if ip[9] != 6:
            continue

        src_ip = '.'.join(str(x) for x in ip[12:16])
        dst_ip = '.'.join(str(x) for x in ip[16:20])
        tcp = ip[ihl:]
        src_port = struct.unpack('>H', tcp[0:2])[0]
        dst_port = struct.unpack('>H', tcp[2:4])[0]
        seq = struct.unpack('>I', tcp[4:8])[0]
        off = ((tcp[12] >> 4) & 0xF) * 4
        payload = tcp[off:]

        if src_port != target_port and dst_port != target_port:
            continue
        if not payload:
            continue

        yield (src_ip, src_port, dst_ip, dst_port, seq, payload)


def source_line(password, mountpoint):
    if not mountpoint.startswith('/'):
        mountpoint = '/' + mountpoint
    return f'SOURCE {password} {mountpoint}\r\nSource-Agent: tcpdump-relay\r\n\r\n'.encode()


def connect_caster(host, port, password, mountpoint):
    s = socket.create_connection((host, port), timeout=10)
    s.sendall(source_line(password, mountpoint))
    s.settimeout(10)
    return s


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--iface', default='any')
    ap.add_argument('--port', type=int, default=12101)
    ap.add_argument('--b-host', required=True)
    ap.add_argument('--b-port', type=int, default=2101)
    ap.add_argument('--b-pass', required=True)
    ap.add_argument('--b-mount', required=True)
    ap.add_argument('--source-ip', default='')
    args = ap.parse_args()

    tcpdump_cmd = [
        'tcpdump', '-i', args.iface, '-s', '0', '-U', '-w', '-',
        f'tcp port {args.port}'
    ]

    print('Start tcpdump:', ' '.join(tcpdump_cmd), file=sys.stderr)
    proc = subprocess.Popen(tcpdump_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    sock = None
    buffers = defaultdict(bytes)
    synced = set()
    last_reconnect = 0

    try:
        for src_ip, src_port, dst_ip, dst_port, seq, payload in parse_pcap_stream(proc.stdout, args.port):
            if args.source_ip and src_ip != args.source_ip and dst_ip != args.source_ip:
                continue

            flow = (src_ip, src_port, dst_ip, dst_port)
            b = buffers[flow] + payload

            if flow not in synced:
                d3 = b.find(b'\xD3')
                if d3 >= 0:
                    b = b[d3:]
                    synced.add(flow)
                else:
                    if len(b) > 8192:
                        b = b[-4096:]
                    buffers[flow] = b
                    continue

            buffers[flow] = b

            if sock is None:
                now = time.time()
                if now - last_reconnect < 1:
                    continue
                last_reconnect = now
                try:
                    sock = connect_caster(args.b_host, args.b_port, args.b_pass, args.b_mount)
                    print('Đã kết nối caster B', file=sys.stderr)
                except Exception as ex:
                    print('Kết nối caster B lỗi:', ex, file=sys.stderr)
                    sock = None
                    continue

            try:
                sock.sendall(b)
                buffers[flow] = b''
            except Exception as ex:
                print('Gửi lỗi, sẽ reconnect:', ex, file=sys.stderr)
                try:
                    sock.close()
                except Exception:
                    pass
                sock = None

    finally:
        try:
            if sock:
                sock.close()
        except Exception:
            pass
        proc.kill()


if __name__ == '__main__':
    main()
```

### 3.4 Chạy demo