# qstate

> :warning: Experimental draft!

This repository contains a qstate encoder and decoder written in Go.
Qstate is a structured state description schema for QUIC connections.
It provides a unified platform for logging, analyzing, serializing, transmitting, and restoring QUIC states.
Qstate complements qlog, where qlog describes the events and state changes over time, qstate represents a snapshot at a particular instant.
Therefore, qstate can give a better overview of the state machine at a glance.
At some point, there may also be tools that generate qstate from qlog at any position.
In addition, qstate is intended to be comprehensive enough to support live connection migration between machines and QUIC implementations.
The design and syntax is inspired by the qlog schema.
Two formats are currently supported: JSON and MessagePack.


Qstate is devided into in three decoupled sections: `transport`, `crypto`, and `metrics`.
Depending on the use case, only parts of the state can be shared.


```json
{
  "transport": {
    "version": 1,
    "chosen_alpn": "perf",
    "vantage_point": "server",
    "connection_ids": [
      {
        "sequence_number": 10,
        "connection_id": "09bdf22f44ae6bc2c12765112ec53a9f482c",
        "stateless_reset_token": "9a0de4ae5e35f02c97d562fafd3d8300"
      },
      {
        "sequence_number": 11,
        "connection_id": "09bdf22f44ae86ac5d900218c9abb3dcc5cf",
        "stateless_reset_token": "9f3f0ee819e4ddaaef6048a496615bec"
      },
      {
        "sequence_number": 9,
        "connection_id": "09bdf22f44aeab2cc14e58703cd6e9e069a1",
        "stateless_reset_token": "6092813a57d7e277484b2cd2a6a9dcaa"
      }
    ],
    "remote_connection_ids": [
      {
        "sequence_number": 5,
        "connection_id": "18b54aaf",
        "stateless_reset_token": "29fbecc5e3767de50a21c5b53a242352"
      },
      {
        "sequence_number": 6,
        "connection_id": "290cb816",
        "stateless_reset_token": "bc5661b5d7f32acd16f5b6cbfcab17cc"
      },
      {
        "sequence_number": 7,
        "connection_id": "e5e5a6bc",
        "stateless_reset_token": "dba11689415c4b0afed9c8574416f524"
      }
    ],
    "dst_ip": "127.0.0.1",
    "dst_port": 48851,
    "parameters": {
      "initial_max_stream_data_bidi_local": 524288,
      "initial_max_stream_data_bidi_remote": 524288,
      "initial_max_stream_data_uni": 524288,
      "max_ack_delay": 26,
      "disable_active_migration": true,
      "max_udp_payload_size": 0,
      "max_idle_timeout": 30000,
      "OriginalDestinationConnectionID": "42c33e2be6eed516f36039838bb5beac",
      "active_connection_id_limit": 4
    },
    "remote_parameters": {
      "initial_max_stream_data_bidi_local": 524288,
      "initial_max_stream_data_bidi_remote": 524288,
      "initial_max_stream_data_uni": 524288,
      "max_ack_delay": 26,
      "disable_active_migration": true,
      "max_udp_payload_size": 1452,
      "max_idle_timeout": 5000,
      "OriginalDestinationConnectionID": null,
      "active_connection_id_limit": 4
    },
    "max_data": 103039344,
    "remote_max_data": 56675274,
    "sent_data": 55958608,
    "received_data": 100000000,
    "max_bidirectional_streams": 397,
    "max_unidirectional_streams": 397,
    "remote_max_bidirectional_streams": 398,
    "remote_max_unidirectional_streams": 398,
    "next_unidirectional_stream": 3,
    "next_bidirectional_stream": 1,
    "remote_next_unidirectional_stream": 2,
    "remote_next_bidirectional_stream": 4,
    "streams": [
      {
        "stream_id": 0,
        "write_offset": 55958608,
        "write_max_data": 56234118,
        "write_ack": 55958607,
        "WriteQueue": null,
        "read_offset": 100000000,
        "read_fin": 100000000,
        "read_max_data": 101731416,
        "ReadQueue": null
      }
    ],
    "next_packet_number": 43322,
    "highest_observed_packet_number": 72465,
    "ack_ranges": [[72459, 72465]],
    "remote_ack_ranges": null,
    "pending_acks": [
      {
        "packet_number": 43312,
        "frames": [
          {
            "frame_type": "max_data"
          }
        ]
      }
    ]
  },
  "crypto": {
    "key_phase": 1,
    "tls_cipher": "TLS_AES_128_GCM_SHA256",
    "remote_header_protection_key": "09c48c7ff217baefc71f54b9ca3dac48",
    "header_protection_key": "f4bf6d026429d30a41cf0387a886d6de",
    "remote_traffic_secret": "69fb76cd9cc063e8afaf8c0b9391305c281ec4b24a81c008819ac0410db155c4",
    "traffic_secret": "82ec934d1c4b18e0e555ca20c90c44b19d52174fca15588317be8e35f77d17ed"
  },
  "metrics": {
    "congestion_window": 416965,
    "smoothed_rtt": 2
  }
}
```