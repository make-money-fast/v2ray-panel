console.log(data);

const green = '#0bbd87';
const gray = '#cccccc'
const red = '#eb5851'

data.editor = null;
data.state =  [{
    content: '服务器运行状态',
    color: gray
}, {
    content: '端口是否启动',
    color: gray
}, {
    content: '本地代理连接状态',
    color: gray
}, {
    content: '能否远程代理',
    color: gray
}];

new Vue({
    el: '#app', data: data, methods: {
        toggleStatus() {
            if (this.isRunning) {
                this.request('/server/api/start', 'get')
            } else {
                this.request('/server/api/stop', 'get')
            }
        }, initDefaultConfig() {
            this.request('/server/api/initDefaultConfig', 'get')
        }, saveConfig() {
            let config = this.editor.getMarkdown()
            this.request('/server/api/config', 'post',{
                config: config,
            })
        }, request(url, method, body, cb) {
            var bodyRaw = null;
            var headers = {};
            if (method.toLowerCase() === 'post') {
                bodyRaw = JSON.stringify(body);
                headers['Content-Type'] = 'application/json'
            }
            return fetch(url, {
                method: method, headers: headers, body: bodyRaw,
            }).then(res => res.json())
                .then(data => {
                    if (data.code === 0 && cb) {
                        cb(data)
                        return
                    }
                    if (data.code === 0) {
                        Swal.fire({
                            title: "提示", text: "操作成功", icon: "success"
                        }).then(() => {
                            window.location.reload()
                        })
                    }
                    if (data.code === -1) {
                        Swal.fire({
                            title: "提示", text: data.msg, icon: "error"
                        })
                    }
                }).catch(err => {
                Swal.fire({
                    title: "错误", text:err, icon: "error"
                })
            })
        },getState() {
            var that = this
            this.request("/server/api/state",'get',null,(res) => {
                if(res.state.isRunning) {
                    this.state[0].color = green;
                } else {
                    this.state[0].color = red;
                }
                if(res.state.isPortListing) {
                    this.state[1].color = green;
                } else {
                    this.state[1].color = red;
                }
                if(res.state.proxyOK) {
                    this.state[2].color = green;
                } else {
                    this.state[2].color = red;
                }
            })
        }
    }, mounted() {
        this.editor = editormd("editor", {
            path: "/static/editor/lib/",
            theme: "dark",
            editorTheme: "pastel-on-dark",
            markdown: this.configJSON,
            height: 740,
            codeFold: true,
            searchReplace: true,
            toolbar: false,             //关闭工具栏
            watch: false,                // 关闭实时预览
            emoji: false,
            taskList: false,
            tocm: false,         // Using [TOCM]
            tex: false,                   // 开启科学公式TeX语言支持，默认关闭
            flowChart: false,             // 开启流程图支持，默认关闭
            sequenceDiagram: false,       // 开启时序/序列图支持，默认关闭,
            imageUpload: false,
        });

        let cp = new ClipboardJS('.copy');

        cp.on('success', function(e) {
            e.clearSelection();
            Swal.fire({
                title: "提示", text: '成功', icon: "success"
            })
        });
    }
})
