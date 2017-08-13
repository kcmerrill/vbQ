workers:
    count: 1
    cmd: |
        touch /tmp/{{ .Name }}
        echo /tmp/{{ .Name }} > email:{{ index .Args "email" }}
