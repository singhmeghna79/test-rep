/*
   Copyright 2016 Vastech SA (PTY) LTD

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package report

const defaultTemplate = `
%use square brackets as golang text templating delimiters
\documentclass{article}
\usepackage{graphicx}
\usepackage[margin=1in]{geometry}
\usepackage{fancyhdr}
\pagestyle{fancy}
\fancyhf{}
\fancyhead[L]{}
\fancyhead[R]{\includegraphics[width=2cm]{../logo.png}}
\fancyfoot[CO]{\small [[.ReportFooter]]}
\fancyfoot[RO]{\makebox[0pt][l]{\hspace{3em}\small{\thepage}}}%

 \setcounter{tocdepth}{4}

\renewcommand{\headrulewidth}{1pt}
\renewcommand{\footrulewidth}{1pt}
\graphicspath{ {images/} }
\begin{document}
\begin{center}
\huge {\vspace*{0.5cm} \textbf{\underline {[[.Title]] Report}}}\\[0.1cm]
\large {\textbf{Report generation date:} \today}\\[1.5cm]
\end{center}

\begin{center}
\begin{tabular}{l  l  }
\textbf{\textit{From date}} & \textbf{\textit{To date}} \\[0.2cm]
\large {[[.FromFormatted]]} & \large{[[.ToFormatted]]}
\end{tabular} \\[0.5cm]
[[range .Panels]][[if .IsSingleStat]]\begin{minipage}{0.3\textwidth}
\includegraphics[width=\textwidth]{image[[.Id]]}
\end{minipage}
[[else]]\par
\vspace{0.5cm}
\includegraphics[width=\textwidth]{image[[.Id]]}
\par
\vspace{0.5cm}
[[end]][[end]]
\end{center}
\end{document}
`
