var app = angular.module('yellow', ['ngRoute', 'ngCookies']);
app.config(['$routeProvider', function($routeProvider){
	$routeProvider.when('/', {
		templateUrl: '/cust/partials/addList.html',
    controller: 'AddCtrl'
	});
}]);

app.controller('AddCtrl', ['$location', '$scope', '$http', function($location, $scope, $http){
$scope.data = {};
$scope.cats = {};
$scope.show = [];
$scope.files = [];
var tmp = $location.absUrl().split("?");
tmp = tmp[1].split("=");
tmp = tmp[1].split("#");
console.log(tmp);
$http.get('/api/getcat').success(function(data, status){
  $scope.cats = data;
});
$scope.add = function(data){
  data.images = $scope.files;
  data.image = $scope.f;
  console.log(data);
	data.Duration = tmp;

  $http.post('/api/addlisting', data).success(function(data, status){
    $scope.data = {};
    $scope.files = [];
    $scope.image = '';
    if(status === 200){
      $location.path('/cust/partials/suc.html');
			//todo move to page to instruct on success and recurring payments
    }
  });
};

$scope.ads = function(data){
data.image = $scope.f;
  $http.post('/api/newAd', data).success(function(data, status){
    $scope.result = data;
    console.log(data);
    if(status === 200){
      //$location.path('/');
    }
  });
};

$scope.change = function(data){e
  console.log("start");
  if(data.Plus === true){
    $scope.show = true;
  }
};

$scope.newfile = function(file){

  var reader = new FileReader();
  reader.onload = function(u){
        //$scope.files.push(u.target.result);
        $scope.$apply(function($scope) {
          $scope.files.push(u.target.result);
          //console.log(u.target.result);
        });
  };
  reader.readAsDataURL(file);

};

$scope.newfile1 = function(file){


  var reader = new FileReader();
  reader.onload = function(u){
        //$scope.files.push(u.target.result);
        $scope.$apply(function($scope) {
          $scope.f = u.target.result;
          //console.log(u.target.result);
        });
  };
  reader.readAsDataURL(file);

};
}]);
