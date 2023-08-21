document.addEventListener('DOMContentLoaded', (event) => {
    var app = new Vue({
        el: '#index',
        data() {
            return {
                posts: null,
                loading: true,
                errored: false,
                page: this.getPage(),
                pageCount: 0,
                searchQuery: this.getQuery()
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
                .get('/api/news/search?page=' + this.page + '&query=' + encodeURIComponent(this.searchQuery))
                .then(response => (
                    this.posts = response.data.posts,
                    this.pageCount = response.data.page_count
                ))
                .catch(error => {
                    console.log(error);
                    this.errored = true;
                })
                .finally(() => (this.loading = false));
            },
            getQuery() {
                query_string = window.location.search.replace(/^\?/, '')
                query_params = query_string.split('&')
                for (i = 0; i < query_params.length; i++) {
                    param_data = query_params[i].split('=')
                    if (param_data.length == 2 && param_data[0] == "query") {
                        return decodeURIComponent(param_data[1]);
                    }
                }
                return "";
            },
            getPage() {
                page = 1;
                params = window.location.hash.replace(/^\#/, '').split('/')
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
                this.searchQuery = this.getQuery()
                this.fetchNews();    
            }
        },
        mounted() {
            this.page = this.getPage()
            this.searchQuery = this.getQuery()
            this.fetchNews();
        }
    });
    window.onhashchange = app.updatePage
})
