events {
  worker_connections 1024;
}
http {

  {{range $rc := .Items}}
    upstream {{$rc.Name}} {
          {{range $points := (index $rc.Subsets 0).Addresses}}
      server {{$points.IP}}:{{(index (index $rc.Subsets 0).Ports 0).Port}};
                {{end}}
    }
  {{end}}
        server {
    listen 80;
    {{range $rc := .Items}}
                location /{{$rc.Name}}/ {
                  proxy_pass http://{{$rc.Name}}/;
          }
                {{end}}
  }

}
