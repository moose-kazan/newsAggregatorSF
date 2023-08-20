document.addEventListener('DOMContentLoaded', (event) => {
    var app = new Vue({
        el: '#full',
        data() {
            return {
                posts: null,
                loading: true,
                errored: false,
                id: this.getId(),
                nocomments: false,
                comments: null,
                newComment: '',
                commentResultSuccess: '',
                commentResultError: '',
            };
        },
        methods: {
            addComment() {
                if (this.newComment == "") {
                    // Do nothing on empty comment
                    return
                }
                this.commentResultSuccess = ""
                this.commentResultError = ""
                const postData = { comment: this.newComment, id: this.id };
                axios
                .post("/api/comments/add", postData)
                .then(response => ( 
                    this.commentResultSuccess = response.data.success ? response.data.message : '',
                    this.commentResultError = response.data.success ? '' : response.data.message
                ))
                .catch(error => (
                    this.commentResultError = error.message
                ))
                .finally(() => (
                    this.newComment = ""
                ))
            },
            fetchData() {
                axios
                .get('/api/news/detail/' + this.id)
                .then(response => (
                    this.posts = [response.data.Post],
                    this.comments = response.data.Comments,
                    this.nocomments = (this.comments == null || this.comments.length == 0)
                ))
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
                this.fetchData();
            }
        },
        mounted() {
            this.fetchData();
        }
    });
    window.onhashchange = app.updatePage
})
