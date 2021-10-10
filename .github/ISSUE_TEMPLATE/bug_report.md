---
name: Bug report
about: Create a report to help us improve
title: 'Bug: '
labels: bug
assignees: ''

---

**Bug description**
A clear and concise description of what the bug is.

**Step by step reproduction**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Minimal example OpenApi specification file to reproduce**
```
openapi: 3.0.0
...
paths:
  /abcde:
    post:
      summary: this is an endpoint
      description: ''
      operationId: abcde
      responses:
        '405':
          description: Invalid input
          content:
            application/json:
              schema: {}
      security:
        - apikey_auth
      requestBody:
        $ref: '#/components/requestBodies/A'
      parameters: []
```

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem. E.g. for issues with the visualization in tools like redoc or swagger editor.
