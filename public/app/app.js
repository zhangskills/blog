angular.module('my-app', ['ui.router', 'restangular'])

.config(['$stateProvider', '$urlRouterProvider', '$locationProvider',
    function($stateProvider, $urlRouterProvider, $locationProvider) {

        $urlRouterProvider.otherwise('/');

        $stateProvider.state('bloglist', {
            url: '/',
            templateUrl: 'public/app/blog_list.html',
            controller: 'blogListController'
        }).state('blogShow', {
            url: '/blog/:id',
            templateUrl: 'public/app/blog_show.html',
            controller: 'blogShowController'
        });

        $locationProvider.html5Mode(true);

    }
])

.controller('blogListController', function($scope, Restangular) {
    window.document.title = '全部 | Jack\'s Blog';
    Restangular.all('api/blog').get('list').then(function(res) {
        console.log(res.val.blogs);
        $scope.blogs = res.val.blogs;
    });
})

.controller('blogShowController', function($scope, $stateParams, Restangular) {
    Restangular.all('api/blog').get('show', {
        id: $stateParams.id
    }).then(function(res) {
        console.log(res.val);
        $scope.blog = res.val;
        window.document.title = $scope.blog.title + ' | Jack\'s Blog';
    });
});