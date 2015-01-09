angular.module('my-app', ['ui.router', 'restangular'])

.config(['$stateProvider', '$urlRouterProvider', '$locationProvider',
    function($stateProvider, $urlRouterProvider, $locationProvider) {
        $urlRouterProvider.otherwise('/');

        $stateProvider.state('bloglist', {
            url: '/',
            templateUrl: 'public/app/blog_list.html',
            controller: 'blogListController'
        }).state('blogshow', {
            url: '/blog/:id',
            templateUrl: 'public/app/blog_show.html',
            controller: 'blogShowController'
        }).state('tagcloud', {
            url: '/tag/cloud',
            templateUrl: 'public/app/tag_cloud.html',
            controller: 'tagCloudController'
        }).state('about', {
            url: '/about',
            templateUrl: 'public/app/about.html',
            controller: 'aboutController'
        });

        $locationProvider.html5Mode(true);

    }
])

.factory('Rest', ['Restangular', function(Restangular) {
    return Restangular.all('api');
}])

.controller('indexController', function($scope, Rest) {
    Rest.get('tag/cloud').then(function(tags) {
        $scope.hotTags = tags.val;
    });
})

.controller('blogListController', function($rootScope, $scope, Rest) {
    window.document.title = '全部 | Jack\'s Blog';
    $rootScope.nav = 'bloglist';
    Rest.get('blog/list').then(function(res) {
        $scope.blogs = res.val.blogs;
    });
})

.controller('blogShowController', function($rootScope, $scope, $stateParams, Rest) {
    $rootScope.nav = 'blogshow';

    Rest.get('blog/show', {
        id: $stateParams.id
    }).then(function(res) {
        $scope.blog = res.val;
        window.document.title = $scope.blog.title + ' | Jack\'s Blog';
    });
    var disqus_shortname = 'zhangskills';
    (function() {
        var dsq = document.createElement('script');
        dsq.type = 'text/javascript';
        dsq.async = true;
        dsq.src = '//' + disqus_shortname + '.disqus.com/embed.js';
        (document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
    })();
})

.controller('tagCloudController', function($rootScope, $scope, Rest) {
    window.document.title = 'Tag Cloud | Jack\'s Blog';
    $rootScope.nav = 'tagcloud';
    Rest.get('tag/cloud').then(function(tags) {
        var words = [];
        for (var i = 0; i < tags.val.length; i++) {
            words.push({
                text: tags.val[i].Key,
                weight: tags.val[i].Count,
                link: '/tag/' + tags.val[i].Key
            })
        }
        jQuery('#tagcloud').jQCloud(words);
    });
})

.controller('aboutController', function($rootScope) {
    window.document.title = 'About me | Jack\'s Blog';
    $rootScope.nav = 'about';
});;