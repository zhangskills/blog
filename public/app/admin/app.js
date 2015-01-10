angular.module('my-app', ['ui.router', 'ui.bootstrap', 'restangular', 'ngSanitize', 'common'])

.config(
    function($stateProvider, $urlRouterProvider) {
        $stateProvider.state('login', {
            url: '/admin/login',
            templateUrl: 'public/app/admin/login.html',
            controller: 'loginController'
        }).state('admin', {
            absract: true,
            url: '/admin',
            templateUrl: 'public/app/admin/index.html'
        }).state('admin.index', {
            url: '/?page',
            templateUrl: 'public/app/admin/blog_list.html',
            controller: 'adminBloglistController'
        });
    }
)

.factory('Rest', ['Restangular', function(Restangular) {
    return Restangular.one('admin/api');
}])

.controller('loginController', function($scope, Rest) {
    window.document.title = '登录页';
    $scope.login = function() {
        $scope.errMsg = '';
        Rest.post('login', {
            name: $scope.userForm.name,
            password: $scope.userForm.password
        }).then(function(res) {
            if (res.err) {
                $scope.errMsg = res.err;
            } else {
                alert('ok');
            }
        });
    };
})

.controller('adminBloglistController', function($scope, Restangular, Rest, $modal, $stateParams) {
    window.document.title = '博客列表 | 管理后台';

    var refreshBlogList = function() {
        Restangular.all('api/blog').get('list', {
            page: $stateParams.page
        }).then(function(res) {
            $scope.blogs = res.val.blogs;
            $scope.$broadcast("pagination", {
                totalItems: res.val.count,
                currentPage: $stateParams.page
            });
        });
    }
    refreshBlogList();

    $scope.showAddBlog = function() {
        $modal.open({
            templateUrl: '/public/app/admin/modal/blog.html',
            controller: 'modalBlogAddController'
        }).result.then(function(msg) {
            refreshBlogList();
        });
    };

    $scope.showUpdateBlog = function(id) {
        $modal.open({
            templateUrl: '/public/app/admin/modal/blog.html',
            controller: 'modalBlogUpdateController',
            resolve: {
                Id: function() {
                    return id;
                }
            }
        }).result.then(function(msg) {
            refreshBlogList();
        });
    };


    $scope.delBlog = function(id) {
        if (confirm('确认要删除该条博客吗？')) {
            Rest.post('blog/del', {
                id: id
            }).then(function(res) {
                if (res.err) {
                    alert(res.err);
                } else {
                    alert('删除成功');
                    refreshBlogList();
                }
            });
        }
    };

})

.controller('modalBlogAddController', function($scope, $modalInstance, Rest, $location) {
    $scope.modaltitle = '添加博客';

    $scope.saveBlog = function() {
        Rest.post('blog/save', {
            blog: {
                Title: $scope.blog.title,
                Content: $scope.blog.content
            },
            tagNames: $scope.tagNames
        }).then(function(res) {
            if (res.err) {
                alert(res.err);
            } else {
                $modalInstance.close('ok');
                $location.url('/admin/');
            }
        });
    };

    $scope.close = function() {
        $modalInstance.dismiss('cancel');
    };

})

.controller('modalBlogUpdateController', function($scope, $modalInstance, Restangular, Rest, Id) {
    $scope.modaltitle = '更新博客';
    $scope.modaltype = 'update';

    Restangular.all('api/blog').get('show', {
        id: Id,
        isadmin: true
    }).then(function(res) {
        $scope.blog = res.val;
        if (res.val.tags) {
            $scope.tagNames = '';
            for (var i = 0; i < res.val.tags.length; i++) {
                if (i > 0) {
                    $scope.tagNames += ','
                }
                $scope.tagNames += res.val.tags[i].name;
            }
        }
    });

    $scope.saveBlog = function() {
        Rest.post('blog/save', {
            blog: {
                Id: $scope.blog.id,
                Title: $scope.blog.title,
                Content: $scope.blog.content
            },
            tagNames: $scope.tagNames
        }).then(function(res) {
            if (res.err) {
                alert(res.err);
            } else {
                $modalInstance.close('cancel');
            }
        });
    };

    $scope.close = function() {
        $modalInstance.dismiss('cancel');
    };

})

.controller('paginationController', function($scope, $log, $location) {
    $scope.itemsPerPage = 10;
    $scope.totalItems = 1;
    $scope.currentPage = 1;

    $scope.pageChanged = function() {
        if ($scope.currentPage > 1) {
            $location.url('/admin/?page=' + $scope.currentPage);
        } else {
            $location.url('/admin/');
        }
    };

    $scope.$on("pagination", function(event, msg) {
        $scope.totalItems = msg.totalItems;
        $scope.currentPage = msg.currentPage;
    });
});