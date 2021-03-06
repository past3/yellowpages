import m from "mithril";
import localforage from 'localforage';
import {UserModel} from '../models/userAuth.js';
import {Analytics} from '../models/analytics.js';

var MenuComponent = {
  view:function(){
    return (
      <div class="">
        <a class="dn pa2 bb b--light-gray hover-bg-light-gray link navy " href="/" oncreate={m.route.link}>Dashboard</a>
        <a class="db pa2 bb b--light-gray hover-bg-light-gray link navy" href="/categories" oncreate={m.route.link}>Categories</a>
        <a class="db pa2 bb b--light-gray hover-bg-light-gray link navy" href="/adverts" oncreate={m.route.link}>Adverts</a>
        <div class=" bb b--light-gray">
          <a class="db pa2  hover-bg-light-gray link navy" href="/listings/new" oncreate={m.route.link}> Add Listing</a>
          <a class="db pa2 hover-bg-light-gray link navy" href="/listings/unapproved" oncreate={m.route.link}>  Unapproved Listings</a>
          <a class="db pa2 hover-bg-light-gray link navy " href="/listings/find" oncreate={m.route.link}> Find Listings</a>
        </div>
      </div>
    );
  }
}

var AdminShell = {
  fixNav:false,
  oncreate:function(vnode){
    Analytics.GetAnalytics()
    var navBar = document.getElementById("fixedNav")
    var navBarOffset = navBar.offsetTop;
    var last_known_scroll_position = 0;
    var ticking = false;
    function CheckPostionAndUpdateNavClass(scroll_pos) {
      // do something with the scroll position
      if (scroll_pos>navBarOffset){
        vnode.state.fixNav = true
        m.redraw()
      }else{
        vnode.state.fixNav = false
        m.redraw()
      }

    }

    window.addEventListener('scroll', function(e) {
      last_known_scroll_position = window.scrollY;
      if (!ticking) {
        window.requestAnimationFrame(function() {
          CheckPostionAndUpdateNavClass(last_known_scroll_position);
          ticking = false;
        });
      }
      ticking = true;
    });

  },
  view:function(vnode){
	return (
	  <section>
  		<section class="  pt3-ns   ph5-ns black-80 bg-yellow">
  		  <div class={"pa2 pv3-ns  w-100  z-5 "+(vnode.state.fixNav===true?"fixed top-0 left-0 bg-yellow shadow-4":"relative-ns")} id="fixedNav">
    			<div class="dib relative">
            <a href="#" class="dib dn-ns black link v-mid mr3  pa2 ba relative" onclick={()=>vnode.state.showNav=!vnode.state.showNav}>☰</a>
              <div class={" right-0 buttom-0 absolute bg-white shadow-m2 pa3 br1 "+(vnode.state.showNav?"db":"dn")}>
                <MenuComponent />
              </div>

    			  <span class="f3 dib v-mid">Calabar<strong>Pages</strong></span>
    			</div>
          <div class="dib v-mid pv2 fr ">
            <div class="dib">
              <a onclick={()=>vnode.state.showProfileNav=!vnode.state.showProfileNav}>
        			  <span class="dib v-mid">{UserModel.User.Username}</span>
                <small class="dib v-mid ph2" style="font-size:8px;">▼</small>
              </a>
              <div class={" right-0 buttom-0 absolute bg-white shadow-m2 pa3 br1 "+(vnode.state.showProfileNav?"db":"dn")}>
                <div class="">
                  <a class="db pa2 bb b--light-gray hover-bg-light-gray link navy pointer" onclick={()=>UserModel.Logout()} >logout</a>
                </div>
              </div>
            </div>
    			</div>
  		  </div>
  		  <div class={"cf pa3 "+(vnode.state.fixNav===true?"pt5":"")}>
    			<div class="tc fr-ns ">
    			  <div class="pa2 pa3-ns shadow-m2 br1 mv3">
    				<div class="tc dib ph3">
    				  <span class="db f3 pa1">
    					{Analytics.Data.UnapprovedListingsCount}
            </span>
    				  <strong class="db">
    					Unapproved
            </strong>
    				</div>
    				<div class="tc dib ph3">
    				  <span class="db f3 pa1">
    					{Analytics.Data.ListingsCount}
            </span>
    				  <strong class="db ">
    					Listings
            </strong>
    				</div>
    				<div class="tc dib ph3">
    				  <span class="db f3 pa1">
    					{Analytics.Data.UsersCount}
            </span>
    				  <strong class="db">
    					Users
            </strong>
    				</div>
    			  </div>
    			</div>
  		  </div>
  		</section>

  		<section  class="ph2 ph5-ns mt-m25-ns pb2 cf" style="min-height:80vh" >
  		  <section class="dib w-30-ns ph3-ns v-top fl">
    			<div class="bg-white shadow-m2 pa3 br1 dn dib-ns w-100">
            <MenuComponent />
    			</div>
        </section><section class="dib w-70-ns ph3-ns v-top br1 fl">
          {m.fragment(vnode.attrs,[vnode.children])}
    		  </section>
  		</section>

      <footer class="pa3 bg-light-gray tc mt4">
        <p class="mb2">Copyright © 2016 Calabar YellowPages </p>
        <small>
          Developed By <a href="http://past3.com.ng" target="_blank">Past3 Software Studios</a> All rights reserved.
        </small>
      </footer>
	  </section>
	)
  }
}
export default AdminShell;
