console.log(data);

const green = '#0bbd87';
const gray = '#cccccc'

data.editor = null;
data.state =  [{
    content: '运行状态',
    color: gray
}, {
    content: '端口是否启动',
    color: gray
}, {
    content: '本地是否能连接代理端口',
    color: gray
}, {
    content: '能否代理远程地址',
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
                switch (res.state) {
                    case 1:
                        that.state[0].color = gray;
                        that.state[1].color = gray;
                        that.state[2].color = gray;
                        that.state[2].color = gray;
                        return
                    case 2:
                        that.state[0].color = green;
                        that.state[1].color = green;
                        that.state[2].color = green;
                        that.state[3].color = gray;
                        return
                    case 3:
                        that.state[0].color = green;
                        that.state[1].color = green;
                        that.state[2].color = green;
                        that.state[3].color = green;
                        return
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

        this.getState()
    }
})
