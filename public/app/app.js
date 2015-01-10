angular.module('my-app', ['ui.router', 'ui.bootstrap', 'restangular', 'ngSanitize'])

.config(['$stateProvider', '$urlRouterProvider', '$locationProvider',
    function($stateProvider, $urlRouterProvider, $locationProvider) {
        $urlRouterProvider.otherwise('/');

        $stateProvider.state('bloglist', {
            url: '/?page',
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
        }).state('tag', {
            url: '/tag/:tag',
            templateUrl: 'public/app/blog_list.html',
            controller: 'tagController'
        }).state('about', {
            url: '/about',
            templateUrl: 'public/app/about.html',
            controller: 'aboutController'
        });

        $locationProvider.html5Mode(true);

    }
])

.run(function($rootScope, $window, $location, $log) {
    $rootScope.$on('$locationChangeSuccess', function(event, url) {
        if (!$location.$$path in ['/', '/tag/cloud', '/about']) {
            $rootScope.nav = '';
        }
        document.documentElement.scrollTop = document.body.scrollTop = 0;
    });
})

.factory('Rest', ['Restangular', function(Restangular) {
    return Restangular.all('api');
}])

.controller('indexController', function($scope, Rest) {
    Rest.get('tag/hot').then(function(tags) {
        $scope.hotTags = tags.val;
    });
    Rest.get('blog/hot').then(function(blogs) {
        $scope.hotBlogs = blogs.val;
    });
})

.controller('blogListController', function($rootScope, $scope, $stateParams, Rest) {
    window.document.title = '全部 | Jack\'s Blog';
    $rootScope.nav = 'bloglist';
    $scope.pageName = '全部';

    Rest.get('blog/list', {
        page: $stateParams.page
    }).then(function(res) {
        $scope.blogs = res.val.blogs;
        $scope.$broadcast("pagination", {
            totalItems: res.val.count,
            currentPage: $stateParams.page
        });
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
    Rest.get('tag/hot').then(function(tags) {
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

.controller('tagController', function($scope, $stateParams, Rest) {
    window.document.title = '标签：' + $stateParams.tag + ' | Jack\'s Blog';
    $scope.pageName = '标签：' + $stateParams.tag;
    Rest.get('tag/blog/list', {
        tagName: $stateParams.tag
    }).then(function(res) {
        $scope.blogs = res.val.blogs;
        $scope.$broadcast("pagination", {
            totalItems: res.val.count,
            currentPage: $stateParams.page
        });
    });
})

.controller('aboutController', function($rootScope) {
    window.document.title = 'About me | Jack\'s Blog';
    $rootScope.nav = 'about';
})

.controller('paginationController', function($scope, $log, $location) {
    $scope.itemsPerPage = 10;
    $scope.totalItems = 1;
    $scope.currentPage = 1;

    $scope.pageChanged = function() {
        if ($scope.currentPage > 1) {
            $location.url('/?page=' + $scope.currentPage);
        } else {
            $location.url('/');
        }
    };

    $scope.$on("pagination", function(event, msg) {
        $scope.totalItems = msg.totalItems;
        $scope.currentPage = msg.currentPage;
    });
});