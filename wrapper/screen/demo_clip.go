package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior/view"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/header"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/clip"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func NewDemoClip() screen.Node {
	title := []line.Line{
		line.TextSpec("=", spec.Cover()),
	}

	clip := clip.New().
		Name("clip - dolor").
		EnableWriteMode().
		SetPause(150).
		SetFrames(
			clip.FrameLines(
				"                                                     ",
				"                   *****++++++****                   ",
				"                 *****++++++++++****                 ",
				"                 ??***++++++++++**??                 ",
				"                 %?****++++++++***?%                 ",
				"                 %%%%%???****???%%%%                 ",
				"                 #SSSSSS%%%%%%SSSSS#                 ",
				"                  %SSSSSSSSSSSSSSS%                  ",
				"                 #%%%%%%%%%%%%%%%%%#                 ",
				"                #%        ##       %#                ",
				"                #?      %?##S?     ?#                ",
				"                #%S##S##S##%#S##S#S%#                ",
				"                  #####SSS%S%SS####                  ",
				"                   ## %?#*  *#?% #                   ",
				"                   SS##*      *##S                   ",
				"                   ##SS#S?##?S#SS#                   ",
				"                      #S%%??%%S#                     ",
				"                                                     ",
			),
			clip.FrameLines(
				"                                                     ",
				"                      **+++++++*                     ",
				"                  ?******++++++**??S                 ",
				"                ???????*****++*****?%S               ",
				"               %?????????*********??%%#S             ",
				"               S???%%%SSS%%???????%%%%SS             ",
				"               S%%%SSSS#SSSSSSSSSS%%%SSS             ",
				"                SSS     ####SSS%%%%SSSS#             ",
				"                 SSSS #   #S####S%*?%#*%             ",
				"                   ##S#  #S%       S%SS              ",
				"                   SS%%?%%%%?**?%SS# #SS             ",
				"                     ###%%%SSSSS  #  ##S             ",
				"                     #####  ##  S%%S%SS              ",
				"                     ####S#                          ",
				"                     ###S### SS?#* %S#               ",
				"                         #############               ",
				"                              #######                ",
				"                                                     ",
			),
			clip.FrameLines(
				"                                                     ",
				"                         ?***??%                     ",
				"                      ??*+++++++***?                 ",
				"                   ???*+++++++++++**???%             ",
				"                  %??*+++++++??***????%%%%S          ",
				"                  %??*++++++++?%????%%%SSSSS         ",
				"                  S%%?*++++++*??%S%SSSSSSS#          ",
				"                  S%%%??***??%%%%S#S#S#####          ",
				"                   SS%%??%?%%%%%%S #####SSS          ",
				"                     S%%%%%%%%S######SS????%         ",
				"                       ###SS#%%?%S##SSS###           ",
				"                            S#S ##SS####             ",
				"                           S %?%%  #%#S              ",
				"                          #%?%%S%     S*             ",
				"                             S##SS%S*                ",
				"                                ####                 ",
				"                                                     ",
				"                                                     ",
			),
			clip.FrameLines(
				"                                                     ",
				"                        ???**????                    ",
				"                     ?*++++++***????S                ",
				"                  ???*++++++++***?*?%%S              ",
				"                 S?**+++++++++***??%SSSS             ",
				"                 %?++++++++++++???SS#%%S             ",
				"                 S%**+++++++*+??%%%%S #S##           ",
				"                 SS%%?*?*+****%%%??%# #S??S          ",
				"                  SSSS%%%%%%%%%S%%%   %#??%          ",
				"                     #SS%SSSSSSSS%SSSSS#             ",
				"                            %#S#%#SSSSS              ",
				"                   ####       %%% S#S#               ",
				"                  S##         %?SS * %               ",
				"                   SS#    ###  S%%%S%*               ",
				"                       S##       S%S                 ",
				"                           S#     SS                 ",
				"                                                     ",
				"                                                     ",
			),
			clip.FrameLines(
				"                                                     ",
				"                      ??????????                     ",
				"                  ????****++****???                  ",
				"                %????++++++++++++???%                ",
				"               %??*+++++++++++++++*??%               ",
				"               ???*+++++++++++++++*???               ",
				"               %*??**++++++++++++*??*%               ",
				"               S*????**++++++++**???*S               ",
				"               #%%???????****??????%%#               ",
				"               #S%%%%????????????%%%S#               ",
				"                SSS%%%%????????%%%SSS                ",
				"                  S    ##SSSS##   S                  ",
				"                   ##S          S#                   ",
				"                   SSS          SS                   ",
				"                   SS%          %S                   ",
				"                   SSS#        #SS                   ",
				"                       SS#####S                      ",
				"                                                     ",
			),
			clip.FrameLines(
				"                                                     ",
				"                     ????**???                       ",
				"                S?????***++++++*?                    ",
				"              S%%?*??***++++++++*??                  ",
				"             SSSS%???***+++++++++*?S                 ",
				"             S%%#SS????+++++++++++?%                 ",
				"            #S# S%%%%%??+*+++++++*%S                 ",
				"           S?S# #%???%%%****+*?*?%SS                 ",
				"           %?#%   %%%%S%%%%%%%%%SSS                  ",
				"             #SSSSS%%SSSSSSSS%SS#                    ",
				"              SSSSS##%#S#%                           ",
				"               #S#S  %%%       ###                   ",
				"               % * SSS?%         #S                  ",
				"               *%S%%%%S  ###    #S                   ",
				"                 S%SS       ##S                      ",
				"                 SS      #S                          ",
				"                                                     ",
				"                                                     ",
			),
			clip.FrameLines(
				"                                                     ",
				"                      %??***?                        ",
				"                 ?****+++++++*??                     ",
				"             %???**++++++++++++*??                   ",
				"           S%%%????****??+++++++*?%                  ",
				"          SSSS%%%?????%?++++++++*?%                  ",
				"           #SSSSSS%SS%??*++++++*?%S                  ",
				"           ####S#S#SS%%%%??***??%%S                  ",
				"           SS##### SS%%%%%%?%??%%S                   ",
				"          %???SS#######S%%%%%%%%S                    ",
				"            ##SSS##SS%?%%#SS###                      ",
				"             ####SS### S#S                           ",
				"              S#%#   %%?% S                          ",
				"             *S      %S%%?%#                         ",
				"                *S%SSS##S                            ",
				"                 #####                               ",
				"                                                     ",
				"                                                     ",
			),
			clip.FrameLines(
				"                                                     ",
				"                      *+++++++**                     ",
				"                 S???**++++++*****?                  ",
				"               S%?******++*****??????                ",
				"             S#%%??**********????????%               ",
				"             SS%%%%????????%%SSS%%???S               ",
				"             SSS%%%SSSSSSSSSSS#SSS%%%S               ",
				"             #SSSS%%%%%SSS####    SSS                ",
				"             %*#%?*%%S####S#   # SSS                 ",
				"              SS%S        %S#  #S#                   ",
				"             SS# #SSS%?**?%%%%?%%S                   ",
				"             S##  #   SSSSS%%%###                    ",
				"              SS%S%%%S  ##  #####                    ",
				"                           #S####                    ",
				"               #S% **#?SS ###S###                    ",
				"               ##############                        ",
				"                ########                             ",
				"                                                     ",
			),
		).
		ToNode()

	clip = view.Map(
		clip,
		mapView,
	)

	return header.Node(clip, title...)
}

func mapView(vm viewmodel.ViewModel) viewmodel.ViewModel {
	kernel := vm.Kernel.ToUnit()

	position := padding.NewBuilder().
		Rows(
			hint.Maximize[winsize.Rows](),
			rows.WithPosition(style.Middle),
		).
		Cols(
			hint.Maximize[winsize.Cols](),
			cols.WithPosition(style.Center),
		).
		ToUnit(kernel)

	vm.Kernel = stack.NewVStack().
		Push(position)

	return vm
}
