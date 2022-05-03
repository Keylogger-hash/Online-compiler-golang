var editor = ace.edit("editor");
let input;
var reader;
editor.setTheme("ace/theme/monokai");
editor.setOptions({
    enableLiveAutocompletion: true,
    enableSnippets: true,
    enableBasicAutocompletion: true
})
editor.session.setMode("ace/mode/golang");
var buttonFormatCode = document.querySelector("#format-code")
var buttonRunCode = document.querySelector("#run-code")
var buttonOpenCode = document.querySelector("#open-code")
var buttonDownloadCode = document.querySelector("#download-code")
window.onbeforeunload = function(){
    text=editor.session.getValue()
    stdoutContent = stdout.textContent
    stderrContent = stderr.textContent
    window.localStorage.setItem('text',text)
    window.localStorage.setItem('stdout',stdoutContent)
    window.localStorage.setItem('stderr',stderrContent)
}
document.addEventListener('DOMContentLoaded',function(){
    text = window.localStorage.getItem('text')
    stdoutContent = window.localStorage.getItem('stdout')
    stderrContent = window.localStorage.getItem('stderr')
    if (text) {
        editor.session.setValue(text)
    }
    if (stdoutContent) {
        stdout.textContent = stdoutContent
    }
    if (stderrContent) {
        stderr.textContent = stderrContent
    }
})
buttonDownloadCode.onclick = function(){
    textValue = editor.getValue();
    file = new Blob([textValue])
    fileLink = document.createElement('a')
    url = URL.createObjectURL(file,{type:'.go'})
    fileLink.href=url;
    fileLink.download='prog.go';
    document.body.appendChild(fileLink);
    fileLink.click()
    setTimeout(function(){
        document.body.removeChild(fileLink) 
        window.URL.revokeObjectURL(url);  
    },0)
}
buttonOpenCode.onclick = function(){
    console.log('open file')
    input=document.createElement('input');
    input.type='file'
    input.accept = ".go"
    input.onchange = _ => {
        let file = Array.from(input.files)[0];
        reader=new FileReader();
        reader.onload = _ =>{
            editor.session.setValue(reader.result);
            console.log('Set session value');
        }
        if (file){
            reader.readAsText(file);
        }
    }
    input.click()
}            
buttonRunCode.onclick = function(){
    var text = editor.session.getValue();
    var stderr = document.querySelectorAll("p.stderr")[0];
    var stdout = document.querySelectorAll("p.stdout")[0];

    var data = {
        Body:text
    }
    var Data = {
        headers: {
            "Content-type":"application/json",
        },
        body:JSON.stringify(data),
        method:"POST"
    }
    stdout.textContent = "Waiting remote server ..."
    fetch("/compile",Data).then(response => response.json()).then((data) => {
        console.log(data)
        if (data["Error"] !== "" ){
            stdout.textContent = data["Error"]
            stderr.textContent = "Go build failed"
        } else {
            stdout.textContent = data["Res"]
            editor.session.setValue(data["Body"],-1)
            stderr.textContent = "Program exited."
        }
    });
}            
buttonFormatCode.onclick = function() {
    var text = editor.getValue();
    var stderr = document.querySelectorAll("p.stderr")[0];
    var stdout = document.querySelectorAll("p.stdout")[0];

    var data = {
        Body:text
    }
    var Data = {
        headers: {
            "Content-type":"application/json",
        },
        body:JSON.stringify(data),
        method:"POST"
    }
    fetch("/fmt",Data).then(response => response.json()).then((data) => {
        console.log(data)
        if (data["Error"] !== "" ){
            stdout.textContent = ""
            stderr.textContent = data["Error"]
            window.localStorage.setItem('stderr',data["Error"])
            window.localStorage.setItem('stdout',"")
        } else {
            stdout.textContent = ""
            stderr.textContent = ""
            window.localStorage.setItem('text',data["Body"])
            window.localStorage.setItem('stderr',"")
            window.localStorage.setItem('stdout',"")
            editor.session.setValue(data["Res"])
        }
    });
}            