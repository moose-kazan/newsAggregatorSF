document.addEventListener('DOMContentLoaded', (event) => {
    var app = new Vue({
        el: '#full',
        data() {
            return {
                posts: null,
                loading: true,
                errored: false,
                id: this.getId()
            };
        },
        methods: {
            fetchNews() {
                axios
                .get('/api/news/detail/' + this.id)
                .then(response => (this.posts = [response.data]))
                .catch(error => {
                    console.log(error);
                    this.errored = true;
                })
                .finally(() => (this.loading = false));
            },
            getId() {
                id = 0;
                params = window.location.hash.replace(/^\#/, '').split('/')
                id = params[0] || 0;
                return id;
            },
            updatePage() {
                this.id = this.getId()
                this.fetchNews();
            }
        },
        mounted() {
            this.fetchNews();
        }
    });
    window.onhashchange = app.updatePage
})
