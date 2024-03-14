console.log(data);

const green = '#0bbd87';
const gray = '#cccccc'
const red = '#eb5851'

data.editor = null;
data.gitHtml = '';
data.url = "https://www.google.com"
data.state =  [{
    content: '运行状态',
    color: gray
},{
    content: 'socks端口',
    color: gray
}, {
    content: 'http端口',
    color: gray
}, {
    content: '能否连上代理服务器',
    color: gray
}, {
    content: '能否通过代理服务器正常上网',
    color: gray
}];

new Vue({
    el: '#app', data: data, methods: {
        toggleStatus() {
            if (this.isRunning) {
                this.request('/client/api/start', 'get')
            } else {
                this.request('/client/api/stop', 'get')
            }
        }, initDefaultConfig() {
            this.request('/client/api/initDefaultConfig', 'get')
        }, saveConfig() {
            let config = this.editor.getMarkdown()
            this.request('/client/api/config', 'post',{
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
            this.request("/client/api/state",'get',null,(res) => {
                if(res.state.isRunning) {
                    this.state[0].color = green
                } else {
                    this.state[0].color = red
                }
                if(res.state.socks) {
                    this.state[1].color = green
                } else {
                    this.state[1].color = red
                }
                if(res.state.http) {
                    this.state[2].color = green
                }else{
                    this.state[2].color = red
                }
                if(res.state.connectToServer) {
                    this.state[3].color = green
                }else{
                    this.state[3].color = red
                }
                if(res.state.porxyOK) {
                    this.state[4].color = green
                }else{
                    this.state[4].color = red
                }
            })
        },
        importVmess() {
            this.request("/client/api/vmess",'post',{vmess: this.vmess})
        },
        checkUrl() {
            this.request("/client/api/check",'post',{url: this.url})
        },configSysProxy() {
            this.request("/client/api/set_proxy?state=on",'get')
        },unsetSysProxy() {
            this.request("/client/api/set_proxy?state=off",'get')
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

        const md = markdownit()
        this.gitHtml = md.render(this.gitTips);
    }
})
