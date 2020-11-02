---
title: "Reva"
date: 2018-05-02T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: _index.md
---

This service provides an ocis extension that wraps [reva](https://github.com/cs3org/reva/) and adds an opinionated configuration to it.

It uses the port range 9140-9179 to preconfigure several services.

| port | service |
+------+---------+
| 9109 | health? |
| 9140 | frontend        |
| 9141 | frontend debug        |
| 9142 | gateway        |
| 9143 | gateway debug        |
| 9144 | users        |
| 9145 | users debug        |
| 9146 | authbasic        |
| 9147 | authbasic debug        |
| 9148 | authbearer        |
| 9149 | authbearer debug        |
| 9150 | sharing        |
| 9151 | sharing debug        |
| 9154 | storage home        |
| 9155 | storage home data        |
| 9156 | storage home debug        |
| 9157 | storage users        |
| 9158 | storage users data        |
| 9159 | storage users debug        |
| 9166-9177 | reserved for s3, wnd, custom + data providers |
| 9178 | storage public link        |
| 9179 | storage public link debug        |