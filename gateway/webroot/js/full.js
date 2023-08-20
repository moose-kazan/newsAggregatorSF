document.addEventListener('DOMContentLoaded', (event) => {
    var app = new Vue({
        el: '#full',
        data() {
            return {
                posts: null,
                loading: true,
                errored: false,
                id: this.getId(),
                commentserrored: false,
                nocomments: false,
                commentsloading: true,
                comments: null,
                newComment: '',
            };
        },
        methods: {
            addComment() {
                if (this.newComment == "") {
                    // Do nothing on empty comment
                    return
                }
                const postData = { comment: this.newComment, id: this.id };
                axios
                .post("/api/comments/add", postData)
                .then(response => ( console.log(response )))
            },
            fetchComments() {
                axios
                .get('/api/comments/last/' + this.id)
                .then(response => (
                    this.comments = response.data,
                    this.nocomments = this.comments ? true : false
                ))
                .catch(error => {
                    console.log(error);
                    this.commentserrored = true;
                })
                .finally(() => (
                    this.commentsloading = false
                ));
            },
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
            this.fetchComments();
        }
    });
    window.onhashchange = app.updatePage
})
