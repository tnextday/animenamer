package kodi

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func Test_UnmarshalTvShow(t *testing.T) {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<tvshow>
    <title>American Gods</title>
    <originaltitle>American Gods</originaltitle>
    <showtitle>American Gods</showtitle>
    <sorttitle>American Gods</sorttitle>
    <ratings>
        <rating name="themoviedb" max="10" default="true">
            <value>6.800000</value>
            <votes>581</votes>
        </rating>
        <rating name="imdb" max="10" default="true">
            <value>5.500000</value>
            <votes>86352</votes>
        </rating>
        <rating name="metacritic" max="10">
            <value>6.0</value>
            <votes>22</votes>
        </rating>
		<rating name="tomatometerallcritics" max="10">
			<value>7.6</value>
			<votes>71</votes>
		</rating>
		<rating name="tomatometerallaudience" max="10">
			<value>6.2</value>
			<votes>119873</votes>
		</rating>
    </ratings>
    <userrating>0</userrating>
    <top250>0</top250>
    <season>2</season>
    <episode>16</episode>
    <displayseason>-1</displayseason>
    <displayepisode>-1</displayepisode>
    <outline></outline>
    <plot>An ex-con becomes the traveling partner of a conman who turns out to be one of the older gods trying to recruit troops to battle the upstart deities. Based on Neil Gaiman&apos;s fantasy novel.</plot>
    <tagline></tagline>
    <runtime>0</runtime>
    <thumb aspect="poster" preview="https://assets.fanart.tv/preview/tv/253573/tvposter/american-gods-58b18cd8d667a.jpg">https://assets.fanart.tv/fanart/tv/253573/tvposter/american-gods-58b18cd8d667a.jpg</thumb>
    <thumb aspect="poster" preview="https://assets.fanart.tv/preview/tv/253573/tvposter/american-gods-5c896dbee9d21.jpg">https://assets.fanart.tv/fanart/tv/253573/tvposter/american-gods-5c896dbee9d21.jpg</thumb>
    <thumb aspect="poster" preview="https://assets.fanart.tv/preview/tv/253573/tvposter/american-gods-57dda913a44e0.jpg">https://assets.fanart.tv/fanart/tv/253573/tvposter/american-gods-57dda913a44e0.jpg</thumb>
    <thumb aspect="poster" preview="https://assets.fanart.tv/preview/tv/253573/tvposter/american-gods-590c159dcbf3a.jpg">https://assets.fanart.tv/fanart/tv/253573/tvposter/american-gods-590c159dcbf3a.jpg</thumb>
    <thumb aspect="banner" preview="https://assets.fanart.tv/preview/tv/253573/tvbanner/american-gods-5cbbdaa84298d.jpg">https://assets.fanart.tv/fanart/tv/253573/tvbanner/american-gods-5cbbdaa84298d.jpg</thumb>
    <thumb aspect="banner" preview="https://assets.fanart.tv/preview/tv/253573/tvbanner/american-gods-5932b1ffb3522.jpg">https://assets.fanart.tv/fanart/tv/253573/tvbanner/american-gods-5932b1ffb3522.jpg</thumb>
    <thumb aspect="banner" preview="https://assets.fanart.tv/preview/tv/253573/tvbanner/american-gods-5932b1ffb43e4.jpg">https://assets.fanart.tv/fanart/tv/253573/tvbanner/american-gods-5932b1ffb43e4.jpg</thumb>
    <thumb aspect="landscape" preview="https://assets.fanart.tv/preview/tv/253573/tvthumb/american-gods-58db45dc886f5.jpg">https://assets.fanart.tv/fanart/tv/253573/tvthumb/american-gods-58db45dc886f5.jpg</thumb>
    <thumb aspect="landscape" preview="https://assets.fanart.tv/preview/tv/253573/tvthumb/american-gods-5932aee79947a.jpg">https://assets.fanart.tv/fanart/tv/253573/tvthumb/american-gods-5932aee79947a.jpg</thumb>
    <thumb aspect="landscape" preview="https://assets.fanart.tv/preview/tv/253573/tvthumb/american-gods-5932aee799e5a.jpg">https://assets.fanart.tv/fanart/tv/253573/tvthumb/american-gods-5932aee799e5a.jpg</thumb>
    <thumb aspect="landscape" preview="https://assets.fanart.tv/preview/tv/253573/tvthumb/american-gods-5932aee79a2f2.jpg">https://assets.fanart.tv/fanart/tv/253573/tvthumb/american-gods-5932aee79a2f2.jpg</thumb>
    <thumb aspect="landscape" preview="https://assets.fanart.tv/preview/tv/253573/tvthumb/american-gods-5932aee79a7c9.jpg">https://assets.fanart.tv/fanart/tv/253573/tvthumb/american-gods-5932aee79a7c9.jpg</thumb>
    <thumb aspect="clearlogo" preview="https://assets.fanart.tv/preview/tv/253573/hdtvlogo/american-gods-58b04bdcecefd.png">https://assets.fanart.tv/fanart/tv/253573/hdtvlogo/american-gods-58b04bdcecefd.png</thumb>
    <thumb aspect="clearlogo" preview="https://assets.fanart.tv/preview/tv/253573/hdtvlogo/american-gods-58b04d78a7ffc.png">https://assets.fanart.tv/fanart/tv/253573/hdtvlogo/american-gods-58b04d78a7ffc.png</thumb>
    <thumb aspect="clearlogo" preview="https://assets.fanart.tv/preview/tv/253573/hdtvlogo/american-gods-59e6660cb7dbc.png">https://assets.fanart.tv/fanart/tv/253573/hdtvlogo/american-gods-59e6660cb7dbc.png</thumb>
    <thumb aspect="clearlogo" preview="https://assets.fanart.tv/preview/tv/253573/hdtvlogo/american-gods-59e6660cc0716.png">https://assets.fanart.tv/fanart/tv/253573/hdtvlogo/american-gods-59e6660cc0716.png</thumb>
    <thumb aspect="clearart" preview="https://assets.fanart.tv/preview/tv/253573/hdclearart/american-gods-59177740ba6cd.png">https://assets.fanart.tv/fanart/tv/253573/hdclearart/american-gods-59177740ba6cd.png</thumb>
    <thumb aspect="clearart" preview="https://assets.fanart.tv/preview/tv/253573/hdclearart/american-gods-5913b6b2ce91d.png">https://assets.fanart.tv/fanart/tv/253573/hdclearart/american-gods-5913b6b2ce91d.png</thumb>
    <thumb aspect="clearart" preview="https://assets.fanart.tv/preview/tv/253573/hdclearart/american-gods-5913b6b2cfa64.png">https://assets.fanart.tv/fanart/tv/253573/hdclearart/american-gods-5913b6b2cfa64.png</thumb>
    <thumb aspect="clearart" preview="https://assets.fanart.tv/preview/tv/253573/hdclearart/american-gods-5913b6b2cf502.png">https://assets.fanart.tv/fanart/tv/253573/hdclearart/american-gods-5913b6b2cf502.png</thumb>
    <thumb aspect="clearart" preview="https://assets.fanart.tv/preview/tv/253573/hdclearart/american-gods-5a4805be0619f.png">https://assets.fanart.tv/fanart/tv/253573/hdclearart/american-gods-5a4805be0619f.png</thumb>
    <thumb aspect="characterart" preview="https://assets.fanart.tv/preview/tv/253573/characterart/american-gods-5a4805af07a04.png">https://assets.fanart.tv/fanart/tv/253573/characterart/american-gods-5a4805af07a04.png</thumb>
    <thumb aspect="characterart" preview="https://assets.fanart.tv/preview/tv/253573/characterart/american-gods-59e6b1c71b65a.png">https://assets.fanart.tv/fanart/tv/253573/characterart/american-gods-59e6b1c71b65a.png</thumb>
    <thumb aspect="poster" type="season" season="2" preview="https://assets.fanart.tv/preview/tv/253573/seasonposter/american-gods-5d1274a8c31cb.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonposter/american-gods-5d1274a8c31cb.jpg</thumb>
    <thumb aspect="poster" type="season" season="1" preview="https://assets.fanart.tv/preview/tv/253573/seasonposter/american-gods-59fea294b565f.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonposter/american-gods-59fea294b565f.jpg</thumb>
    <thumb aspect="poster" type="season" season="1" preview="https://assets.fanart.tv/preview/tv/253573/seasonposter/american-gods-5cacdf37068db.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonposter/american-gods-5cacdf37068db.jpg</thumb>
    <thumb aspect="poster" type="season" season="2" preview="https://assets.fanart.tv/preview/tv/253573/seasonposter/american-gods-5cacdf7783e04.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonposter/american-gods-5cacdf7783e04.jpg</thumb>
    <thumb aspect="poster" type="season" season="2" preview="https://assets.fanart.tv/preview/tv/253573/seasonposter/american-gods-5d1274a8c31cb.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonposter/american-gods-5d1274a8c31cb.jpg</thumb>
    <thumb aspect="poster" type="season" season="1" preview="https://assets.fanart.tv/preview/tv/253573/seasonposter/american-gods-59fea294b565f.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonposter/american-gods-59fea294b565f.jpg</thumb>
    <thumb aspect="poster" type="season" season="1" preview="https://assets.fanart.tv/preview/tv/253573/seasonposter/american-gods-5cacdf37068db.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonposter/american-gods-5cacdf37068db.jpg</thumb>
    <thumb aspect="poster" type="season" season="2" preview="https://assets.fanart.tv/preview/tv/253573/seasonposter/american-gods-5cacdf7783e04.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonposter/american-gods-5cacdf7783e04.jpg</thumb>
    <thumb aspect="banner" type="season" season="1" preview="https://assets.fanart.tv/preview/tv/253573/seasonbanner/american-gods-5cc6b35699d26.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonbanner/american-gods-5cc6b35699d26.jpg</thumb>
    <thumb aspect="banner" type="season" season="2" preview="https://assets.fanart.tv/preview/tv/253573/seasonbanner/american-gods-5cc6b36965b54.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonbanner/american-gods-5cc6b36965b54.jpg</thumb>
    <thumb aspect="banner" type="season" season="1" preview="https://assets.fanart.tv/preview/tv/253573/seasonbanner/american-gods-5cc6b35699d26.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonbanner/american-gods-5cc6b35699d26.jpg</thumb>
    <thumb aspect="banner" type="season" season="2" preview="https://assets.fanart.tv/preview/tv/253573/seasonbanner/american-gods-5cc6b36965b54.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonbanner/american-gods-5cc6b36965b54.jpg</thumb>
    <thumb aspect="landscape" type="season" season="2" preview="https://assets.fanart.tv/preview/tv/253573/seasonthumb/american-gods-5cc6b380d6c56.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonthumb/american-gods-5cc6b380d6c56.jpg</thumb>
    <thumb aspect="landscape" type="season" season="1" preview="https://assets.fanart.tv/preview/tv/253573/seasonthumb/american-gods-59e6b5a03e7aa.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonthumb/american-gods-59e6b5a03e7aa.jpg</thumb>
    <thumb aspect="landscape" type="season" season="2" preview="https://assets.fanart.tv/preview/tv/253573/seasonthumb/american-gods-5cc6b380d6c56.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonthumb/american-gods-5cc6b380d6c56.jpg</thumb>
    <thumb aspect="landscape" type="season" season="1" preview="https://assets.fanart.tv/preview/tv/253573/seasonthumb/american-gods-59e6b5a03e7aa.jpg">https://assets.fanart.tv/fanart/tv/253573/seasonthumb/american-gods-59e6b5a03e7aa.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/m6qf6lq3yARgbZwspvDLbUFtASh.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/gevw5nZRYz2kWj1PqW9pz4sgeeZ.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/btwTe5cQbGWGOErBiRqnjNP9cJl.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/loJ4sfr4zp995qMoeCHiIIGaOg8.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/dHo8Lw7ruIaQTdTTDZPCMyZxwy5.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/zfAXP4bG2G17VuLNU9cqRcVU0xj.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/oxYUbNpG2st2zXWzYRvewehmvuj.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/mwoQ6zynu2DBxKCBYi30qoM236N.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/8XEoXAMzgcf7m1KiUDZ9N1UGh4o.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/rWsayJB1grML2LdPjjKDC3g0Brr.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/8qRsj8uJ4zPARQmQ9FvejTY1lnV.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/acjnZP0GrwWDxCxV6QejKizbzOy.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/hN1sI57QILGfdrEOqpUfo0NtHjW.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/hz2jNy3DfseYzRSybGRlUtz4pTi.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/hLDgNDdrkB0oWiuClpxN4E3XadJ.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/4FiqawHsVz1mYCRudPtXKbfmP4M.jpg</thumb>
    <thumb aspect="poster">http://image.tmdb.org/t/p/original/sKR8Q36YBtyRc19y4yGYuD1xBgA.jpg</thumb>
    <thumb aspect="poster" type="season" season="2">http://image.tmdb.org/t/p/original/4l8Vnbb7e5QA6bAItMqQIHXLRgc.jpg</thumb>
    <thumb aspect="poster" type="season" season="2">http://image.tmdb.org/t/p/original/ni0thXw5Zi5dQKBY6Oj0vcfIS2n.jpg</thumb>
    <thumb aspect="poster" type="season" season="2">http://image.tmdb.org/t/p/original/v17HfCzWKQKOBrww9RxZmN5R9tF.jpg</thumb>
    <thumb aspect="poster" type="season" season="2">http://image.tmdb.org/t/p/original/2ffvlgYsxbXGiWkc3V6Q8tgpiBo.jpg</thumb>
    <thumb aspect="poster" type="season" season="1">http://image.tmdb.org/t/p/original/rASj7OUjWDhfhAeO2MaFOA3lJpQ.jpg</thumb>
    <thumb aspect="poster" type="season" season="1">http://image.tmdb.org/t/p/original/67exRijfvN5RRmBCqFtk1bhJ7Uh.jpg</thumb>
    <thumb aspect="poster" type="season" season="1">http://image.tmdb.org/t/p/original/59iE3xxP7H8rAiXW6TDR2HSoUUm.jpg</thumb>
    <thumb aspect="poster" type="season" season="2">http://image.tmdb.org/t/p/original/4l8Vnbb7e5QA6bAItMqQIHXLRgc.jpg</thumb>
    <thumb aspect="poster" type="season" season="2">http://image.tmdb.org/t/p/original/ni0thXw5Zi5dQKBY6Oj0vcfIS2n.jpg</thumb>
    <thumb aspect="poster" type="season" season="2">http://image.tmdb.org/t/p/original/v17HfCzWKQKOBrww9RxZmN5R9tF.jpg</thumb>
    <thumb aspect="poster" type="season" season="2">http://image.tmdb.org/t/p/original/2ffvlgYsxbXGiWkc3V6Q8tgpiBo.jpg</thumb>
    <thumb aspect="banner">https://thetvdb.com/banners/graphical/253573-g3.jpg</thumb>
    <thumb aspect="banner">https://thetvdb.com/banners/graphical/253573-g4.jpg</thumb>
    <thumb aspect="banner">https://thetvdb.com/banners/graphical/253573-g2.jpg</thumb>
    <thumb aspect="banner">https://thetvdb.com/banners/graphical/253573-g.jpg</thumb>
    <thumb aspect="banner">https://thetvdb.com/banners/graphical/253573-g5.jpg</thumb>
    <fanart>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-5c8965c58e778.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-5c8965c58e778.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-59e6a8a495c2a.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-59e6a8a495c2a.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-59e6b13827ba2.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-59e6b13827ba2.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-5932b089e07ad.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-5932b089e07ad.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-5932b089e2913.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-5932b089e2913.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-5932b089e0000.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-5932b089e0000.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-5932b089e0d3a.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-5932b089e0d3a.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-5932b089e1395.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-5932b089e1395.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-5932b089e1952.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-5932b089e1952.jpg</thumb>
        <thumb preview="https://assets.fanart.tv/preview/tv/253573/showbackground/american-gods-5932b089e23ca.jpg">https://assets.fanart.tv/fanart/tv/253573/showbackground/american-gods-5932b089e23ca.jpg</thumb>
    </fanart>
    <mpaa>Australia:MA</mpaa>
    <playcount>0</playcount>
    <lastplayed></lastplayed>
    <episodeguide>
        <url cache="tmdb-46639-en.json">http://api.themoviedb.org/3/tv/46639?api_key=6a5be4999abf74eba1f9a8311294c267&amp;language=en</url>
    </episodeguide>
    <id>46639</id>
    <uniqueid type="tmdb" default="true">46639</uniqueid>
    <uniqueid type="tvdb">253573</uniqueid>
    <genre>Drama</genre>
    <genre>Mystery</genre>
    <genre>Sci-Fi &amp; Fantasy</genre>
    <premiered>2017-04-30</premiered>
    <year>2017</year>
    <status></status>
    <code></code>
    <aired></aired>
    <studio>Starz</studio>
    <trailer></trailer>
    <actor>
        <name>Ricky Whittle</name>
        <role>Shadow Moon</role>
        <order>0</order>
        <thumb>http://image.tmdb.org/t/p/original/cjeDbVfBp6Qvb3C74Dfy7BKDTQN.jpg</thumb>
    </actor>
    <actor>
        <name>Ian McShane</name>
        <role>Mr. Wednesday</role>
        <order>1</order>
        <thumb>http://image.tmdb.org/t/p/original/pY9ud4BJwHekNiO4MMItPbgkdAy.jpg</thumb>
    </actor>
    <actor>
        <name>Emily Browning</name>
        <role>Laura Moon</role>
        <order>2</order>
        <thumb>http://image.tmdb.org/t/p/original/fa1Kyj02wxwcdS6EHb2i27TNXvU.jpg</thumb>
    </actor>
    <actor>
        <name>Pablo Schreiber</name>
        <role>Mad Sweeney</role>
        <order>3</order>
        <thumb>http://image.tmdb.org/t/p/original/uo8YljeePz3pbj7gvWXdB4gOOW4.jpg</thumb>
    </actor>
    <actor>
        <name>Bruce Langley</name>
        <role>Technical Boy</role>
        <order>4</order>
        <thumb>http://image.tmdb.org/t/p/original/f4EOWUmznLqboq8Ce7jnlkHVK3Y.jpg</thumb>
    </actor>
    <actor>
        <name>Yetide Badaki</name>
        <role>Bilquis</role>
        <order>5</order>
        <thumb>http://image.tmdb.org/t/p/original/qfzkREHuI1JvMxBteIAjKX8qMEr.jpg</thumb>
    </actor>
    <namedseason number="1">Season 1</namedseason>
    <namedseason number="2">Season 2</namedseason>
    <resume>
        <position>0.000000</position>
        <total>0.000000</total>
    </resume>
    <dateadded>2017-10-07 14:25:47</dateadded>
</tvshow>`
	var tvshow TVShow
	err := xml.Unmarshal([]byte(data), &tvshow)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	fmt.Printf("%+v", tvshow)
	//for _, r := range tvshow.Ratings {
	//	fmt.Println(r)
	//}
	//for _, f := range tvshow.Fanarts {
	//	fmt.Println(f)
	//}
}

func Test_UnmarshalEpisodeDetails(t *testing.T) {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<episodedetails>
    <title>The Bone Orchard</title>
    <showtitle>American Gods</showtitle>
    <ratings>
        <rating name="tmdb" max="10" default="true">
            <value>7.532000</value>
            <votes>31</votes>
        </rating>
    </ratings>
    <userrating>0</userrating>
    <top250>0</top250>
    <season>1</season>
    <episode>1</episode>
    <displayseason>-1</displayseason>
    <displayepisode>-1</displayepisode>
    <outline></outline>
    <plot>When Shadow Moon is released from prison early after the death of his wife, he meets Mr. Wednesday and is recruited as his bodyguard. Shadow discovers that this may be more than he bargained for.</plot>
    <tagline></tagline>
    <runtime>0</runtime>
    <thumb>http://image.tmdb.org/t/p/original/uvry4weK00pFLn7fxQ9M4m3Da2A.jpg</thumb>
    <mpaa>16</mpaa>
    <playcount>0</playcount>
    <lastplayed></lastplayed>
    <id>1276153</id>
    <uniqueid type="tmdb" default="true">1276153</uniqueid>
    <genre>Drama</genre>
    <genre>Mystery</genre>
    <genre>Sci-Fi &amp; Fantasy</genre>
    <credits>Bryan Fuller</credits>
    <credits>Michael Green</credits>
    <director>David Slade</director>
    <premiered>2017-04-30</premiered>
    <year>2017</year>
    <status></status>
    <code></code>
    <aired>2017-04-30</aired>
    <studio>Starz</studio>
    <trailer></trailer>
    <actor>
        <name>Jonathan Tucker</name>
        <role>&apos;Low Key&apos; Lyesmith</role>
        <order>10</order>
        <thumb>http://image.tmdb.org/t/p/original/jvJpYDbwmUTACw7Yn7PKOP6CdlJ.jpg</thumb>
    </actor>
    <actor>
        <name>Demore Barnes</name>
        <role>Mr. Ibis</role>
        <order>11</order>
        <thumb>http://image.tmdb.org/t/p/original/4rEVzSIFPgiN14xYQnjKcKQ7tYE.jpg</thumb>
    </actor>
    <actor>
        <name>Betty Gilpin</name>
        <role>Audrey</role>
        <order>12</order>
        <thumb>http://image.tmdb.org/t/p/original/xFeqyem5i4Kf0nFjBZ4Oi9NM26k.jpg</thumb>
    </actor>
    <actor>
        <name>Beth Grant</name>
        <role>Jack</role>
        <order>13</order>
        <thumb>http://image.tmdb.org/t/p/original/zAT9GvzJE0ytL3C36L461cgKI9p.jpg</thumb>
    </actor>
    <actor>
        <name>Joel Murray</name>
        <role>Paunch</role>
        <order>14</order>
        <thumb>http://image.tmdb.org/t/p/original/t5syYfCgxbTC7XPrNeXhhhQULUf.jpg</thumb>
    </actor>
    <actor>
        <name>Ricky Whittle</name>
        <role>Shadow Moon</role>
        <order>0</order>
        <thumb>http://image.tmdb.org/t/p/original/cjeDbVfBp6Qvb3C74Dfy7BKDTQN.jpg</thumb>
    </actor>
    <actor>
        <name>Ian McShane</name>
        <role>Mr. Wednesday</role>
        <order>1</order>
        <thumb>http://image.tmdb.org/t/p/original/pY9ud4BJwHekNiO4MMItPbgkdAy.jpg</thumb>
    </actor>
    <actor>
        <name>Emily Browning</name>
        <role>Laura Moon</role>
        <order>2</order>
        <thumb>http://image.tmdb.org/t/p/original/fa1Kyj02wxwcdS6EHb2i27TNXvU.jpg</thumb>
    </actor>
    <actor>
        <name>Pablo Schreiber</name>
        <role>Mad Sweeney</role>
        <order>3</order>
        <thumb>http://image.tmdb.org/t/p/original/uo8YljeePz3pbj7gvWXdB4gOOW4.jpg</thumb>
    </actor>
    <actor>
        <name>Bruce Langley</name>
        <role>Technical Boy</role>
        <order>4</order>
        <thumb>http://image.tmdb.org/t/p/original/f4EOWUmznLqboq8Ce7jnlkHVK3Y.jpg</thumb>
    </actor>
    <actor>
        <name>Yetide Badaki</name>
        <role>Bilquis</role>
        <order>5</order>
        <thumb>http://image.tmdb.org/t/p/original/qfzkREHuI1JvMxBteIAjKX8qMEr.jpg</thumb>
    </actor>
    <resume>
        <position>0.000000</position>
        <total>0.000000</total>
    </resume>
    <dateadded>2017-10-07 14:25:47</dateadded>
</episodedetails>`
	var episodedetails EpisodeDetails
	err := xml.Unmarshal([]byte(data), &episodedetails)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	fmt.Printf("%+v", episodedetails)
}
