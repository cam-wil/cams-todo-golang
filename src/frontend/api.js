function getAll() {
  fetch('http://127.0.0.1:10000/todos')
  .then(function (response) {
    return response.json();
  })
  .then(function (d) {
    let col = ["id", "name", "content", "complete"]
    let dat = "<table class=\"table\"><tr><th scope=\"col\">id</th><th scope=\"col\">name</th><th scope=\"col\">content</th><th scope=\"col\">complete</th><th scope=\"col\">Delete</th></tr><tr>"
    let div = document.getElementById("dvTable")

    for(let i = 0; i < d.length; i++) {
      dat += "<tr scope=\"row\">"
      for(let j = 0; j < 4; j++) {
        let c = col[j]
        dat += `<td>${d[i][c]}</td>`
        //console.log(d)
      }
      dat += `<td><a onclick=\"deleteTodo(${d[i]["id"]})\"href=\"\">delete</a> </td></tr>`
    }
    dat += "</tr></table>"
    div.innerHTML = ""
    div.innerHTML = dat
  })
  .catch(function (err) {
    console.log(err);
  });
}
// Name, Content, Complete
function addTodo(name, content, complete) {
  const dat = JSON.stringify({id: 0, name: name, content: content, complete: complete})
  fetch(`http://127.0.0.1:10000/todo`, {
    method: 'POST',
    mode: 'no-cors', 
    body: dat
  }).then(d => {
    getAll()
    document.getElementById('name').value = ""
    document.getElementById('content').value = ""
    document.getElementById('complete').checked = false
  })
}

function completeTodo(id) {
  alert("complete " + id)
}

function deleteTodo(id) {
  const url = `http://127.0.0.1:10000/todos/d/${id}`
  fetch(url, {
    method: 'POST',
    headers: {
      'Access-Control-Allow-Origin':'*'
  },
    mode: 'no-cors',
  }).then((r) => {
    alert('delete')
    console.log(r)
  }).catch(err => console.error(err))
 
}