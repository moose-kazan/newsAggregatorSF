document.addEventListener('DOMContentLoaded', (event) => {
    var app = new Vue({
        el: '#index',
        data() {
            return {
                posts: null,
                loading: true,
                errored: false,
                page: this.getPage()
            };
        },
        methods: {
            pageNext() {
                if (this.posts.length > 0) {
                    this.page++;
                    window.location.hash = "page=" + this.page;
                }
            },
            pagePrev() {
                if (this.page > 1) {
                    this.page--;
                    window.location.hash = "page=" + this.page;
                }
            },
            fetchNews() {
                axios
                .get('/api/news/latest?page=' + this.page)
                .then(response => (this.posts = response.data))
                .catch(error => {
                    console.log(error);
                    this.errored = true;
                })
                .finally(() => (this.loading = false));
            },
            getPage() {
                page = 1;
                params = window.location.hash.replace(/^\#/, '').split('/')
                console.log(params)
                if (params.length < 1) {
                    return page;
                }
                for (i = 0; i < params.length; i++) {
                    paramName = params[i].split('=')[0]
                    if (paramName != 'page') {
                        continue
                    }
                    page = params[i].split('=')[1] || 1;
                    break
                }
                return page;
            },
            updatePage() {
                this.page = this.getPage()
                this.fetchNews();    
            }
        },
        mounted() {
            this.fetchNews();
        }
    });
    window.onhashchange = app.updatePage
})
